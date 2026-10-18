package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/skybedy/laravel-tene.life/internal/i18n"
	"github.com/skybedy/laravel-tene.life/internal/models"
	"github.com/skybedy/laravel-tene.life/internal/newsai"
	"github.com/skybedy/laravel-tene.life/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type fakeNews struct {
	article        models.NewsArticle
	saved, deleted int
	public         bool
	offset         int
}

type fakeArticleAI struct {
	draft newsai.Draft
	input newsai.Input
	calls int
}

func (f *fakeArticleAI) Enabled() bool { return true }
func (f *fakeArticleAI) Model() string { return "test-model" }
func (f *fakeArticleAI) Generate(_ context.Context, input newsai.Input) (newsai.Draft, error) {
	f.calls++
	f.input = input
	return f.draft, nil
}

func (f *fakeNews) Categories(context.Context) ([]models.NewsCategory, error) {
	return []models.NewsCategory{{ID: 1, Name: "Doprava", Slug: "doprava"}}, nil
}
func (f *fakeNews) List(_ context.Context, _ string, public bool, _ int, offset int) ([]models.NewsArticle, int, error) {
	f.public = public
	f.offset = offset
	return []models.NewsArticle{f.article}, 25, nil
}
func (f *fakeNews) BySlug(_ context.Context, slug string) (models.NewsArticle, error) {
	if f.article.Status != "published" || slug != f.article.Slug {
		return models.NewsArticle{}, sql.ErrNoRows
	}
	return f.article, nil
}
func (f *fakeNews) ByID(context.Context, int64) (models.NewsArticle, error) { return f.article, nil }
func (f *fakeNews) Save(_ context.Context, a *models.NewsArticle) error {
	f.saved++
	a.ID = 1
	f.article = *a
	return nil
}
func (f *fakeNews) Delete(context.Context, int64, int) error { f.deleted++; return nil }
func (f *fakeNews) Sitemap(context.Context) ([]store.NewsSitemapEntry, error) {
	if f.article.Status != "published" {
		return nil, nil
	}
	return []store.NewsSitemapEntry{{Slug: f.article.Slug, UpdatedAt: f.article.UpdatedAt}}, nil
}

func newsTestServer(t *testing.T) (*echo.Echo, *fakeNews) {
	return newsTestServerWithAI(t, nil)
}

func newsTestServerWithAI(t *testing.T, ai ArticleAI) (*echo.Echo, *fakeNews) {
	t.Helper()
	e := echo.New()
	templates := template.New("").Funcs(template.FuncMap{"localeURL": i18n.LocaleURL, "languageFlag": i18n.LanguageFlag})
	templates, err := templates.ParseGlob("../../views/news-*.html")
	if err != nil {
		t.Fatal(err)
	}
	templates, err = templates.ParseFiles("../../views/header.html", "../../views/nav.html", "../../views/footer.html")
	if err != nil {
		t.Fatal(err)
	}
	e.Renderer = &TemplateRenderer{Templates: templates}
	now := time.Now().UTC()
	f := &fakeNews{article: models.NewsArticle{ID: 1, CategoryID: 1, Title: "Zpráva <script>alert(1)</script>", Slug: "zprava", Summary: "Popis & události", Body: "První odstavec\n\n<script>alert(1)</script>", Status: "published", PublishedAt: &now, UpdatedAt: now, Category: models.NewsCategory{ID: 1, Name: "Doprava", Slug: "doprava"}, ImageURL: "/images/test.jpg", ImageAlt: "Popis snímku"}}
	hash, err := bcrypt.GenerateFromPassword([]byte("test-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	h := &NewsHandler{Store: f, AI: ai, Config: NewsConfig{BaseURL: "http://localhost", Username: "editor", PasswordHash: string(hash)}}
	h.Register(e)
	return e, f
}

func TestNewsPublicRenderingAndSEO(t *testing.T) {
	e, f := newsTestServer(t)
	for _, path := range []string{"/zpravy", "/zpravy?kategorie=doprava&page=2", "/zpravy/zprava", "/sitemap.xml", "/robots.txt"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if rec.Code != 200 {
			t.Fatalf("%s: %d %s", path, rec.Code, rec.Body)
		}
		body := rec.Body.String()
		if strings.Contains(body, "<script>alert(1)</script>") {
			t.Fatalf("XSS in %s", path)
		}
		if path == "/zpravy?kategorie=doprava&page=2" && (f.offset != 12 || !f.public || !strings.Contains(body, "kategorie=doprava")) {
			t.Fatal("pagination/filter lost")
		}
		if path == "/zpravy/zprava" {
			for _, expected := range []string{`rel="canonical" href="http://localhost/zpravy/zprava"`, `property="og:title"`, `property="og:description"`, `content="http://localhost/images/test.jpg"`, `application/ld+json`} {
				if !strings.Contains(body, expected) {
					t.Fatalf("missing %s", expected)
				}
			}
			start := strings.Index(body, `<script type="application/ld+json">`) + len(`<script type="application/ld+json">`)
			end := strings.Index(body[start:], "</script>") + start
			var ld map[string]any
			if err := json.Unmarshal([]byte(body[start:end]), &ld); err != nil || ld["@type"] != "NewsArticle" {
				t.Fatalf("JSON-LD: %v %v", ld, err)
			}
			if strings.Contains(body, "<script>alert(1)</script>") || !strings.Contains(body, "<!-- raw HTML omitted -->") {
				t.Fatal("unsafe Markdown HTML was rendered")
			}
		}
	}
	for _, path := range []string{"/zpravy?kategorie=nezname", "/zpravy?page=0", "/zpravy?page=999999999999999999", "/zpravy/neexistuje"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if rec.Code < 400 {
			t.Fatalf("invalid path accepted: %s", path)
		}
	}
	f.article.Status = "draft"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest("GET", "/zpravy/zprava", nil))
	if rec.Code != 404 {
		t.Fatal("draft exposed")
	}
}

func TestNewsAIDraftFillsFormWithoutSaving(t *testing.T) {
	ai := &fakeArticleAI{draft: newsai.Draft{
		Title: "Dopravní omezení na Tenerife", Summary: "Silnici čeká krátká uzavírka.",
		Body:         "## Kdy omezení platí\n\nSilnice bude uzavřena.\n\n- první bod\n- druhý bod",
		CategorySlug: "doprava", SEOTitle: "Uzavírka silnice na Tenerife",
		SEODescription: "Termín a dopady plánované uzavírky.", EditorNotes: "Ověřit přesný čas.",
	}}
	e, store := newsTestServerWithAI(t, ai)
	get := httptest.NewRequest("GET", "/admin/zpravy/novy", nil)
	get.SetBasicAuth("editor", "test-password")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, get)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "AI asistent článku") {
		t.Fatalf("AI editor unavailable: %d %s", rec.Code, rec.Body)
	}
	var csrf *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "news_csrf" {
			csrf = cookie
		}
	}
	if csrf == nil {
		t.Fatal("missing CSRF cookie")
	}
	form := url.Values{
		"csrf": {csrf.Value}, "revision": {"0"}, "status": {"published"},
		"ai_source": {"Zdrojové podklady k uzavírce."}, "ai_instructions": {"Zdůrazni objížďku."},
		"action": {"ai_generate"},
	}
	req := httptest.NewRequest("POST", "/admin/zpravy/novy", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("editor", "test-password")
	req.AddCookie(csrf)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	body := rec.Body.String()
	if rec.Code != http.StatusOK || ai.calls != 1 || store.saved != 0 {
		t.Fatalf("AI generation: status=%d calls=%d saves=%d body=%s", rec.Code, ai.calls, store.saved, body)
	}
	for _, expected := range []string{"Dopravní omezení na Tenerife", "dopravni-omezeni-na-tenerife", "Kdy omezení platí", "Ověřit přesný čas", `value="draft" selected`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("AI form missing %q", expected)
		}
	}
	if ai.input.SourceMaterial != "Zdrojové podklady k uzavírce." || ai.input.EditorialBrief != "Zdůrazni objížďku." {
		t.Fatalf("unexpected AI input: %+v", ai.input)
	}
}

func TestNewsMarkdownRendererBlocksRawHTMLAndDangerousLinks(t *testing.T) {
	rendered, err := renderNewsMarkdown("## Nadpis\n\n**Tučně** a [bezpečně](https://tene.life).\n\n<script>alert(1)</script>\n\n[útok](javascript:alert(2))")
	if err != nil {
		t.Fatal(err)
	}
	html := string(rendered)
	for _, expected := range []string{"<h2>Nadpis</h2>", "<strong>Tučně</strong>", `href="https://tene.life"`, "<!-- raw HTML omitted -->"} {
		if !strings.Contains(html, expected) {
			t.Fatalf("missing %q in %s", expected, html)
		}
	}
	if strings.Contains(html, "<script") || strings.Contains(html, `href="javascript:`) {
		t.Fatalf("unsafe Markdown output: %s", html)
	}
}

func TestNewsAdminAuthenticationCSRFAndPreview(t *testing.T) {
	e, f := newsTestServer(t)
	for _, path := range []string{"/admin/zpravy", "/admin/zpravy/novy", "/admin/zpravy/1", "/admin/zpravy/1/smazat"} {
		for _, method := range []string{"GET", "POST"} {
			if path == "/admin/zpravy" && method == "POST" {
				continue
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
			if rec.Code != 401 && rec.Code != 405 {
				t.Fatalf("unprotected %s %s: %d", method, path, rec.Code)
			}
		}
	}
	get := httptest.NewRequest("GET", "/admin/zpravy/novy", nil)
	get.SetBasicAuth("editor", "test-password")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, get)
	if rec.Code != 200 {
		t.Fatalf("editor: %d %s", rec.Code, rec.Body)
	}
	var csrf *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "news_csrf" {
			csrf = cookie
		}
	}
	if csrf == nil {
		t.Fatal("missing csrf cookie")
	}
	if !csrf.HttpOnly || csrf.SameSite != http.SameSiteStrictMode {
		t.Fatal("cookie flags")
	}
	form := url.Values{"title": {"Zpráva z přírody"}, "summary": {"Perex"}, "body": {"Text"}, "category_id": {"1"}, "status": {"published"}, "action": {"preview"}}
	post := func(path string, form url.Values, site string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if site != "" {
			req.Header.Set("Sec-Fetch-Site", site)
		}
		req.SetBasicAuth("editor", "test-password")
		req.AddCookie(csrf)
		rr := httptest.NewRecorder()
		e.ServeHTTP(rr, req)
		return rr
	}
	if rr := post("/admin/zpravy/novy", form, ""); rr.Code != 400 && rr.Code != 403 {
		t.Fatalf("missing CSRF accepted: %d", rr.Code)
	}
	form.Set("csrf", csrf.Value)
	if rr := post("/admin/zpravy/novy", form, "cross-site"); rr.Code != 403 {
		t.Fatalf("cross-site accepted: %d", rr.Code)
	}
	rr := post("/admin/zpravy/novy", form, "")
	if rr.Code != 200 || f.saved != 0 || !strings.Contains(rr.Body.String(), "Náhled") || strings.Contains(rr.Body.String(), "application/ld+json") {
		t.Fatalf("preview: %d saves=%d %s", rr.Code, f.saved, rr.Body)
	}
	if rr.Header().Get("Cache-Control") != "private, no-store" || rr.Header().Get("X-Robots-Tag") == "" {
		t.Fatal("preview caching/indexing")
	}
	form.Set("action", "save")
	rr = post("/admin/zpravy/novy", form, "")
	if rr.Code != 303 || f.saved != 1 || f.article.Slug != "zprava-z-prirody" {
		t.Fatalf("save: %d %s %+v", rr.Code, rr.Body, f.article)
	}
	rr = post("/admin/zpravy/1/smazat", form, "")
	if rr.Code != 400 || f.deleted != 0 {
		t.Fatal("unconfirmed deletion")
	}
	form.Set("confirm", "yes")
	form.Set("revision", "1")
	rr = post("/admin/zpravy/1/smazat", form, "")
	if rr.Code != 303 || f.deleted != 1 {
		t.Fatalf("confirmed deletion: %d", rr.Code)
	}
}

func TestNewsValidation(t *testing.T) {
	if got := newsSlug("  Příliš žluťoučký kůň – Tenerife 2026! "); got != "prilis-zlutoucky-kun-tenerife-2026" {
		t.Fatal(got)
	}
	for _, raw := range []string{"javascript:alert(1)", "data:image/svg+xml,test", "//evil.test/a.jpg", "/images/../private", "/images/%2e%2e/private", "https://user:pass@example.com/a", "https://example.com\\evil"} {
		if validNewsURL(raw, true) {
			t.Fatalf("unsafe image: %s", raw)
		}
	}
	for _, raw := range []string{"https://example.com/image.jpg", "/images/news/test.jpg"} {
		if !validNewsURL(raw, true) {
			t.Fatalf("valid image rejected: %s", raw)
		}
	}
	a := models.NewsArticle{Title: "Titulek", Slug: "titulek", Summary: "P", Body: "Text", CategoryID: 1, Status: "draft"}
	cats := []models.NewsCategory{{ID: 1}}
	if err := validateNews(a, cats); err != "" {
		t.Fatal(err)
	}
	a.SourceURL = "https://example.com"
	if validateNews(a, cats) == "" {
		t.Fatal("source pair missing")
	}
	a.SourceURL = ""
	a.CategoryID = 99
	if validateNews(a, cats) == "" {
		t.Fatal("invalid category")
	}
	a.CategoryID = 1
	a.Body = strings.Repeat("ž", 50001)
	if validateNews(a, cats) == "" {
		t.Fatal("body length")
	}
}

func TestNewsAdminDisabledAndRateLimited(t *testing.T) {
	e := echo.New()
	h := &NewsHandler{}
	h.Register(e)
	rr := httptest.NewRecorder()
	e.ServeHTTP(rr, httptest.NewRequest("GET", "/admin/zpravy", nil))
	if rr.Code != 404 {
		t.Fatal("unconfigured admin enabled")
	}
	e, _ = newsTestServer(t)
	for i := 0; i < 21; i++ {
		rr = httptest.NewRecorder()
		e.ServeHTTP(rr, httptest.NewRequest("GET", "/admin/zpravy", nil))
	}
	if rr.Code != 429 {
		t.Fatal("auth rate limiting failed")
	}
}

func TestNewsConfiguration(t *testing.T) {
	t.Setenv("NEWS_ADMIN_USERNAME", "")
	t.Setenv("NEWS_ADMIN_PASSWORD_HASH", "")
	for _, origin := range []string{"", "https://tene.life", "http://localhost:8080", "http://127.0.0.1:8080"} {
		t.Setenv("SITE_URL", origin)
		if _, err := NewsConfigFromEnv(); err != nil {
			t.Fatalf("valid origin %q: %v", origin, err)
		}
	}
	for _, origin := range []string{"http://tene.life", "https://tene.life/path", "https://user:secret@tene.life", "javascript:alert(1)", "https://tene.life?query=1"} {
		t.Setenv("SITE_URL", origin)
		if _, err := NewsConfigFromEnv(); err == nil {
			t.Fatalf("invalid origin accepted: %q", origin)
		}
	}
	t.Setenv("SITE_URL", "https://tene.life")
	t.Setenv("NEWS_ADMIN_USERNAME", "editor")
	if _, err := NewsConfigFromEnv(); err == nil {
		t.Fatal("missing hash accepted")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("example-password"), 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("NEWS_ADMIN_PASSWORD_HASH", string(hash))
	if _, err := NewsConfigFromEnv(); err != nil {
		t.Fatal(err)
	}
}
