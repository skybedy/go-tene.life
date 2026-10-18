package web

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skybedy/laravel-tene.life/internal/i18n"
	"github.com/skybedy/laravel-tene.life/internal/models"
	"github.com/skybedy/laravel-tene.life/internal/newsai"
	"github.com/skybedy/laravel-tene.life/internal/store"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/time/rate"
)

type NewsRepository interface {
	Categories(context.Context) ([]models.NewsCategory, error)
	List(context.Context, string, bool, int, int) ([]models.NewsArticle, int, error)
	BySlug(context.Context, string) (models.NewsArticle, error)
	ByID(context.Context, int64) (models.NewsArticle, error)
	Save(context.Context, *models.NewsArticle) error
	Delete(context.Context, int64, int) error
	Sitemap(context.Context) ([]store.NewsSitemapEntry, error)
}

type ArticleAI interface {
	Enabled() bool
	Model() string
	Generate(context.Context, newsai.Input) (newsai.Draft, error)
}

type NewsConfig struct{ BaseURL, Username, PasswordHash string }

func NewsConfigFromEnv() (NewsConfig, error) {
	c := NewsConfig{strings.TrimRight(strings.TrimSpace(os.Getenv("SITE_URL")), "/"), os.Getenv("NEWS_ADMIN_USERNAME"), os.Getenv("NEWS_ADMIN_PASSWORD_HASH")}
	if c.BaseURL == "" {
		c.BaseURL = "https://tene.life"
	}
	u, err := url.Parse(c.BaseURL)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
		return c, errors.New("SITE_URL must be an absolute origin without path/query")
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" || u.Hostname() == "::1")) {
		return c, errors.New("SITE_URL requires HTTPS (except local loopback)")
	}
	if c.Username != "" || c.PasswordHash != "" {
		cost, err := bcrypt.Cost([]byte(c.PasswordHash))
		if c.Username == "" || err != nil || cost < 10 || cost > 14 {
			return c, errors.New("news admin requires username and bcrypt hash with cost 10–14")
		}
	}
	return c, nil
}

type NewsHandler struct {
	Store     NewsRepository
	AI        ArticleAI
	Config    NewsConfig
	aiLimiter *rate.Limiter
}

type newsPage struct {
	models.PageData
	Articles                                              []models.NewsArticle
	Categories                                            []models.NewsCategory
	Article                                               models.NewsArticle
	BodyHTML                                              template.HTML
	CategorySlug, Canonical, Description, OGImage, OGType string
	JSONLD                                                template.JS
	CSRF, Error, FormAction, PreviousURL, NextURL         string
	AISource, AIInstructions, AIMessage, AIModel          string
	Page, Total                                           int
	Preview, Saved, AIEnabled                             bool
}

func (h *NewsHandler) page(c echo.Context, title string, admin bool) newsPage {
	section := "news"
	if admin {
		section = "news-admin"
	}
	d := newsPage{PageData: models.PageData{PageTitle: title, Locale: "cs", CurrentPath: c.Request().URL.Path, CurrentSection: section, I18n: i18n.Messages("cs")}, OGType: "website"}
	if !admin {
		d.GAMeasurementID = strings.TrimSpace(os.Getenv("GA_MEASUREMENT_ID"))
		d.GAEnabled = d.GAMeasurementID != ""
	}
	if admin && h.AI != nil && h.AI.Enabled() {
		d.AIEnabled = true
		d.AIModel = h.AI.Model()
	}
	d.CSRF, _ = c.Get("csrf").(string)
	return d
}

func (h *NewsHandler) Register(e *echo.Echo) {
	if h.AI != nil && h.AI.Enabled() {
		h.aiLimiter = rate.NewLimiter(rate.Every(15*time.Second), 2)
	}
	// Revalidation prevents shared caches from continuing to expose withdrawn stories.
	public := e.Group("", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error { c.Response().Header().Set("Cache-Control", "no-cache"); return next(c) }
	})
	public.GET("/zpravy", h.List)
	public.GET("/zpravy/:slug", h.Detail)
	public.GET("/sitemap.xml", h.Sitemap)
	e.GET("/robots.txt", func(c echo.Context) error {
		return c.String(200, "User-agent: *\nDisallow: /admin/\nSitemap: "+h.Config.BaseURL+"/sitemap.xml\n")
	})
	admin := e.Group("/admin/zpravy", h.adminHeaders, h.authentication(), middleware.BodyLimit("512K"), middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf", CookieName: "news_csrf", CookiePath: "/admin/zpravy", CookieHTTPOnly: true, CookieSecure: strings.HasPrefix(h.Config.BaseURL, "https://"), CookieSameSite: http.SameSiteStrictMode,
	}))
	admin.GET("", h.AdminList)
	admin.GET("/novy", h.Edit)
	admin.POST("/novy", h.Save)
	admin.GET("/:id", h.Edit)
	admin.POST("/:id", h.Save)
	admin.GET("/:id/smazat", h.ConfirmDelete)
	admin.POST("/:id/smazat", h.Delete)
}

func (h *NewsHandler) adminHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "private, no-store")
		c.Response().Header().Set("X-Robots-Tag", "noindex, nofollow")
		c.Response().Header().Set("Referrer-Policy", "no-referrer")
		c.Response().Header().Set("X-Frame-Options", "DENY")
		return next(c)
	}
}

func (h *NewsHandler) authentication() echo.MiddlewareFunc {
	// One bounded bucket for the single editorial account. No spoofable proxy/IP keys.
	limiter := rate.NewLimiter(rate.Every(2*time.Second), 20)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if h.Config.Username == "" || h.Config.PasswordHash == "" {
				return echo.ErrNotFound
			}
			if !limiter.Allow() {
				c.Response().Header().Set("Retry-After", "2")
				return echo.NewHTTPError(429, "Příliš mnoho pokusů. Zkuste to za chvíli.")
			}
			user, pass, ok := c.Request().BasicAuth()
			a, b := sha256.Sum256([]byte(user)), sha256.Sum256([]byte(h.Config.Username))
			validPassword := false
			if ok && len(pass) <= 72 {
				validPassword = bcrypt.CompareHashAndPassword([]byte(h.Config.PasswordHash), []byte(pass)) == nil
			}
			if !ok || subtle.ConstantTimeCompare(a[:], b[:]) != 1 || !validPassword {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="Tenelife zpravy", charset="UTF-8"`)
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func newsPageNumber(c echo.Context) (int, error) {
	s := c.QueryParam("page")
	if s == "" {
		return 1, nil
	}
	p, err := strconv.Atoi(s)
	if err != nil || p < 1 || p > 100000 {
		return 0, echo.NewHTTPError(400, "Neplatná stránka.")
	}
	return p, nil
}

func newsPagination(d *newsPage, base string, page, total, size int, category string) {
	d.Page = page
	d.Total = total
	link := func(p int) string {
		q := url.Values{"page": {strconv.Itoa(p)}}
		if category != "" {
			q.Set("kategorie", category)
		}
		return base + "?" + q.Encode()
	}
	if page > 1 {
		d.PreviousURL = link(page - 1)
	}
	if page*size < total {
		d.NextURL = link(page + 1)
	}
}

func (h *NewsHandler) List(c echo.Context) error {
	d := h.page(c, "Zprávy z Tenerife", false)
	p, err := newsPageNumber(c)
	if err != nil {
		return err
	}
	d.Categories, err = h.Store.Categories(c.Request().Context())
	if err != nil {
		return err
	}
	d.CategorySlug = c.QueryParam("kategorie")
	if d.CategorySlug != "" {
		found := false
		for _, cat := range d.Categories {
			if cat.Slug == d.CategorySlug {
				found = true
				d.PageTitle = cat.Name + " – Zprávy z Tenerife"
			}
		}
		if !found {
			return echo.ErrNotFound
		}
	}
	articles, total, err := h.Store.List(c.Request().Context(), d.CategorySlug, true, 12, (p-1)*12)
	if err != nil {
		return err
	}
	if p > 1 && len(articles) == 0 {
		return echo.ErrNotFound
	}
	d.Articles = articles
	newsPagination(&d, "/zpravy", p, total, 12, d.CategorySlug)
	q := url.Values{}
	if p > 1 {
		q.Set("page", strconv.Itoa(p))
		d.PageTitle += fmt.Sprintf(" – strana %d", p)
	}
	if d.CategorySlug != "" {
		q.Set("kategorie", d.CategorySlug)
	}
	d.Canonical = h.Config.BaseURL + "/zpravy"
	if len(q) > 0 {
		d.Canonical += "?" + q.Encode()
	}
	d.Description = "Zprávy z Tenerife: doprava, počasí a příroda, život na ostrově, akce, kultura a turistika."
	return c.Render(200, "news-list.html", d)
}

func (h *NewsHandler) articleData(d *newsPage, a models.NewsArticle, public bool) {
	d.Article = a
	d.BodyHTML, _ = renderNewsMarkdown(a.Body)
	d.PageTitle = a.SEOTitle
	if d.PageTitle == "" {
		d.PageTitle = a.Title
	}
	d.Description = a.SEODescription
	if d.Description == "" {
		d.Description = a.Summary
	}
	if !public {
		return
	}
	d.Canonical = h.Config.BaseURL + "/zpravy/" + a.Slug
	d.OGImage = a.ImageURL
	d.OGType = "article"
	if strings.HasPrefix(d.OGImage, "/") {
		d.OGImage = h.Config.BaseURL + d.OGImage
	}
	data := map[string]any{"@context": "https://schema.org", "@type": "NewsArticle", "headline": a.Title, "description": d.Description, "datePublished": a.PublishedISO(), "dateModified": a.UpdatedAt.UTC().Format(time.RFC3339), "mainEntityOfPage": d.Canonical, "url": d.Canonical, "inLanguage": "cs", "articleSection": a.Category.Name, "publisher": map[string]string{"@type": "Organization", "name": "Tenelife", "url": h.Config.BaseURL}}
	if d.OGImage != "" {
		data["image"] = []string{d.OGImage}
	}
	if a.SourceURL != "" {
		data["isBasedOn"] = a.SourceURL
	}
	// encoding/json escapes HTML delimiters, including </script> from untrusted text.
	encoded, _ := json.Marshal(data)
	d.JSONLD = template.JS(encoded)
}

func (h *NewsHandler) Detail(c echo.Context) error {
	a, err := h.Store.BySlug(c.Request().Context(), c.Param("slug"))
	if errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	d := h.page(c, a.Title, false)
	h.articleData(&d, a, true)
	recent, _, err := h.Store.List(c.Request().Context(), "", true, 4, 0)
	if err != nil {
		return err
	}
	for _, other := range recent {
		if other.ID != a.ID && len(d.Articles) < 3 {
			d.Articles = append(d.Articles, other)
		}
	}
	return c.Render(200, "news-detail.html", d)
}

func (h *NewsHandler) AdminList(c echo.Context) error {
	d := h.page(c, "Správa zpráv", true)
	p, err := newsPageNumber(c)
	if err != nil {
		return err
	}
	a, total, err := h.Store.List(c.Request().Context(), "", false, 20, (p-1)*20)
	if err != nil {
		return err
	}
	if p > 1 && len(a) == 0 {
		return echo.ErrNotFound
	}
	d.Articles = a
	newsPagination(&d, "/admin/zpravy", p, total, 20, "")
	return c.Render(200, "news-admin-list.html", d)
}

func (h *NewsHandler) editData(c echo.Context) (newsPage, error) {
	d := h.page(c, "Upravit zprávu", true)
	d.FormAction = c.Request().URL.Path
	var err error
	d.Categories, err = h.Store.Categories(c.Request().Context())
	if err != nil {
		return d, err
	}
	d.Article.Status = "draft"
	if c.Param("id") != "" {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || id < 1 {
			return d, echo.ErrNotFound
		}
		d.Article, err = h.Store.ByID(c.Request().Context(), id)
		if errors.Is(err, sql.ErrNoRows) {
			return d, echo.ErrNotFound
		}
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func (h *NewsHandler) Edit(c echo.Context) error {
	d, err := h.editData(c)
	if err != nil {
		return err
	}
	d.Saved = c.QueryParam("ulozeno") == "1"
	return c.Render(200, "news-edit.html", d)
}

var newsSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func newsSlug(title string) string {
	var b strings.Builder
	dash := false
	for _, r := range norm.NFD.String(strings.ToLower(title)) {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			if dash && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			dash = false
		} else {
			dash = true
		}
	}
	s := b.String()
	if len(s) > 180 {
		s = strings.TrimRight(s[:180], "-")
	}
	return s
}

func validNewsURL(raw string, image bool) bool {
	if raw == "" {
		return true
	}
	if strings.ContainsAny(raw, "\\\r\n\t") {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil || u.User != nil || u.Fragment != "" {
		return false
	}
	if image && u.Scheme == "" && u.Host == "" && strings.HasPrefix(raw, "/images/") && !strings.Contains(u.Path, "..") {
		return true
	}
	return u.Hostname() != "" && (u.Scheme == "https" || !image && u.Scheme == "http")
}

func validateNews(a models.NewsArticle, categories []models.NewsCategory) string {
	for _, f := range []struct {
		name, value string
		max         int
		required    bool
	}{
		{"Titulek", a.Title, 200, true}, {"Slug", a.Slug, 180, true}, {"Perex", a.Summary, 500, true}, {"Obsah", a.Body, 50000, true},
		{"Obrázek", a.ImageURL, 2048, false}, {"Popis obrázku", a.ImageAlt, 200, false}, {"Název zdroje", a.SourceName, 200, false}, {"URL zdroje", a.SourceURL, 2048, false}, {"SEO titulek", a.SEOTitle, 200, false}, {"SEO popis", a.SEODescription, 500, false},
	} {
		if !utf8.ValidString(f.value) || strings.ContainsRune(f.value, 0) || utf8.RuneCountInString(f.value) > f.max || f.required && strings.TrimSpace(f.value) == "" {
			return fmt.Sprintf("%s: vyplňte platný text (nejvýše %d znaků).", f.name, f.max)
		}
	}
	if !newsSlugPattern.MatchString(a.Slug) {
		return "Slug smí obsahovat malá písmena a–z, číslice a spojovníky."
	}
	if a.Status != "draft" && a.Status != "published" {
		return "Neplatný stav článku."
	}
	found := false
	for _, cat := range categories {
		if cat.ID == a.CategoryID {
			found = true
		}
	}
	if !found {
		return "Vyberte platnou kategorii."
	}
	if !validNewsURL(a.ImageURL, true) {
		return "Obrázek musí mít HTTPS URL nebo cestu /images/…"
	}
	if a.ImageURL != "" && a.ImageAlt == "" {
		return "Doplňte popis obrázku."
	}
	if !validNewsURL(a.SourceURL, false) || (a.SourceName == "") != (a.SourceURL == "") {
		return "U zdroje vyplňte název i platnou HTTP(S) URL."
	}
	return ""
}

func (h *NewsHandler) Save(c echo.Context) error {
	d, err := h.editData(c)
	if err != nil {
		return err
	}
	f := func(key string) string { return strings.TrimSpace(c.FormValue(key)) }
	d.AISource = strings.TrimSpace(c.FormValue("ai_source"))
	d.AIInstructions = strings.TrimSpace(c.FormValue("ai_instructions"))
	a := d.Article
	a.Title = f("title")
	a.Summary = f("summary")
	a.Body = f("body")
	a.Status = f("status")
	if a.PublishedAt == nil {
		a.Slug = f("slug")
		if a.Slug == "" {
			a.Slug = newsSlug(a.Title)
		}
	}
	a.CategoryID, _ = strconv.ParseInt(f("category_id"), 10, 64)
	a.Revision, _ = strconv.Atoi(f("revision"))
	a.ImageURL = f("image_url")
	a.ImageAlt = f("image_alt")
	a.SourceName = f("source_name")
	a.SourceURL = f("source_url")
	a.SEOTitle = f("seo_title")
	a.SEODescription = f("seo_description")
	d.Article = a
	action := f("action")
	if action == "ai_generate" {
		return h.generateAIDraft(c, &d)
	}
	d.Error = validateNews(a, d.Categories)
	if d.Error != "" {
		return c.Render(422, "news-edit.html", d)
	}
	if action == "preview" {
		for _, cat := range d.Categories {
			if cat.ID == a.CategoryID {
				a.Category = cat
			}
		}
		h.articleData(&d, a, false)
		d.Preview = true
		return c.Render(200, "news-edit.html", d)
	}
	if action != "save" {
		return echo.ErrBadRequest
	}
	if err = h.Store.Save(c.Request().Context(), &a); err != nil {
		if errors.Is(err, store.ErrNewsConflict) {
			d.Error = "Slug už existuje nebo článek mezitím upravil někdo jiný. Zvolte jiný slug, případně znovu načtěte článek."
			return c.Render(409, "news-edit.html", d)
		}
		if errors.Is(err, store.ErrNewsCategory) {
			d.Error = "Kategorie již neexistuje. Vyberte jinou."
			return c.Render(422, "news-edit.html", d)
		}
		return err
	}
	return c.Redirect(303, fmt.Sprintf("/admin/zpravy/%d?ulozeno=1", a.ID))
}

func (h *NewsHandler) generateAIDraft(c echo.Context, d *newsPage) error {
	if h.AI == nil || !h.AI.Enabled() {
		d.Error = "AI asistent není nakonfigurovaný. Doplňte OPENAI_API_KEY a restartujte aplikaci."
		return c.Render(http.StatusServiceUnavailable, "news-edit.html", *d)
	}
	if h.aiLimiter == nil || !h.aiLimiter.Allow() {
		c.Response().Header().Set("Retry-After", "15")
		d.Error = "AI návrh lze spustit nejvýše jednou za 15 sekund. Chvíli počkejte a zkuste to znovu."
		return c.Render(http.StatusTooManyRequests, "news-edit.html", *d)
	}
	if d.AISource == "" || utf8.RuneCountInString(d.AISource) > 40000 {
		d.Error = "Do AI podkladů vložte text o délce nejvýše 40 000 znaků."
		return c.Render(http.StatusUnprocessableEntity, "news-edit.html", *d)
	}
	if utf8.RuneCountInString(d.AIInstructions) > 2000 {
		d.Error = "Pokyny pro AI mohou mít nejvýše 2 000 znaků."
		return c.Render(http.StatusUnprocessableEntity, "news-edit.html", *d)
	}

	categorySlugs := make([]string, 0, len(d.Categories))
	for _, category := range d.Categories {
		categorySlugs = append(categorySlugs, category.Slug)
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 100*time.Second)
	defer cancel()
	draft, err := h.AI.Generate(ctx, newsai.Input{
		SourceMaterial:  d.AISource,
		EditorialBrief:  d.AIInstructions,
		CategoryChoices: categorySlugs,
	})
	if err != nil {
		log.Printf("news AI generation failed: %v", err)
		d.Error = "AI návrh se nepodařilo vytvořit. Podklady zůstaly ve formuláři; zkuste to prosím znovu."
		return c.Render(http.StatusBadGateway, "news-edit.html", *d)
	}

	a := d.Article
	a.Title = strings.TrimSpace(draft.Title)
	a.Summary = strings.TrimSpace(draft.Summary)
	a.Body = strings.TrimSpace(draft.Body)
	a.SEOTitle = strings.TrimSpace(draft.SEOTitle)
	a.SEODescription = strings.TrimSpace(draft.SEODescription)
	if a.PublishedAt == nil {
		a.Slug = newsSlug(a.Title)
	}
	for _, category := range d.Categories {
		if category.Slug == draft.CategorySlug {
			a.CategoryID = category.ID
			a.Category = category
			break
		}
	}
	if a.ID == 0 {
		a.Status = "draft"
	}
	d.Article = a
	if validationError := validateNews(a, d.Categories); validationError != "" {
		log.Printf("news AI generated draft failed local validation: %s", validationError)
		d.Error = "AI vrátila neúplný návrh. Podklady zůstaly ve formuláři; zkuste upravit pokyny a generování zopakovat."
		return c.Render(http.StatusBadGateway, "news-edit.html", *d)
	}
	d.AIMessage = "AI návrh byl vložen do formuláře a zatím není uložený. Pečlivě zkontrolujte fakta, data, názvy a zdroj."
	if strings.TrimSpace(draft.EditorNotes) != "" {
		d.AIMessage += " Poznámka AI: " + strings.TrimSpace(draft.EditorNotes)
	}
	return c.Render(http.StatusOK, "news-edit.html", *d)
}

func (h *NewsHandler) ConfirmDelete(c echo.Context) error {
	d, err := h.editData(c)
	if err != nil {
		return err
	}
	return c.Render(200, "news-delete.html", d)
}

func (h *NewsHandler) Delete(c echo.Context) error {
	if c.FormValue("confirm") != "yes" {
		return echo.NewHTTPError(400, "Potvrďte smazání článku.")
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return echo.ErrNotFound
	}
	revision, err := strconv.Atoi(c.FormValue("revision"))
	if err != nil {
		return echo.ErrBadRequest
	}
	err = h.Store.Delete(c.Request().Context(), id, revision)
	if errors.Is(err, store.ErrNewsConflict) {
		return echo.NewHTTPError(409, "Článek se změnil. Znovu otevřete potvrzení smazání.")
	}
	if err != nil {
		return err
	}
	return c.Redirect(303, "/admin/zpravy")
}

func (h *NewsHandler) Sitemap(c echo.Context) error {
	entries, err := h.Store.Sitemap(c.Request().Context())
	if err != nil {
		return err
	}
	type item struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod,omitempty"`
	}
	doc := struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []item   `xml:"url"`
	}{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, path := range []string{"/", "/zpravy", "/statistics/daily", "/statistics/weekly", "/statistics/monthly", "/statistics/monthly-daily", "/statistics/recent", "/statistics/annual", "/atlantic-sea-temperature-tenerife", "/atlantic-sea-temperature-canary-islands", "/spanelsko-ceska-slovicka"} {
		doc.URLs = append(doc.URLs, item{Loc: h.Config.BaseURL + path})
	}
	for _, e := range entries {
		doc.URLs = append(doc.URLs, item{Loc: h.Config.BaseURL + "/zpravy/" + e.Slug, LastMod: e.UpdatedAt.UTC().Format(time.RFC3339)})
	}
	return c.XML(200, doc)
}
