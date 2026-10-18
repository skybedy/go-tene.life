package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/skybedy/laravel-tene.life/internal/models"
)

// NEWS_TEST_DSN points to a disposable test server. Each run creates its own database.
func TestNewsStoreIntegration(t *testing.T) {
	dsn := os.Getenv("NEWS_TEST_DSN")
	if dsn == "" {
		t.Skip("set NEWS_TEST_DSN for MariaDB integration test")
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.DBName = ""
	cfg.ParseTime = true
	cfg.MultiStatements = true
	server, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	name := fmt.Sprintf("tenelife_news_test_%d", time.Now().UnixNano())
	if _, err = server.Exec("CREATE DATABASE " + name + " CHARACTER SET utf8mb4"); err != nil {
		t.Fatal(err)
	}
	defer server.Exec("DROP DATABASE " + name)
	cfg.DBName = name
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	migration, err := os.ReadFile("../../migrations/0012_news.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	s := &NewsStore{DB: db}
	cats, err := s.Categories(ctx)
	if err != nil || len(cats) != 5 {
		t.Fatalf("categories: %v %v", cats, err)
	}
	a := models.NewsArticle{Title: "První", Slug: "prvni", Summary: "Perex", Body: "Obsah", CategoryID: cats[0].ID, Status: "draft"}
	if err = s.Save(ctx, &a); err != nil {
		t.Fatal(err)
	}
	if _, err = s.BySlug(ctx, a.Slug); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("draft leaked: %v", err)
	}
	if rows, n, err := s.List(ctx, "", true, 12, 0); err != nil || n != 0 || len(rows) != 0 {
		t.Fatalf("draft list: %v %v %v", rows, n, err)
	}
	if rows, err := s.Sitemap(ctx); err != nil || len(rows) != 0 {
		t.Fatalf("draft sitemap: %v %v", rows, err)
	}
	a, err = s.ByID(ctx, a.ID)
	if err != nil {
		t.Fatal(err)
	}
	a.Status = "published"
	if err = s.Save(ctx, &a); err != nil {
		t.Fatal(err)
	}
	a, err = s.BySlug(ctx, a.Slug)
	if err != nil || a.PublishedAt == nil {
		t.Fatalf("publish: %v %v", a, err)
	}
	firstPublished := *a.PublishedAt
	if rows, err := s.Sitemap(ctx); err != nil || len(rows) != 1 {
		t.Fatalf("published sitemap: %v %v", rows, err)
	}
	if rows, n, err := s.List(ctx, cats[1].Slug, true, 12, 0); err != nil || n != 0 || len(rows) != 0 {
		t.Fatalf("category filter: %v %v %v", rows, n, err)
	}
	duplicate := a
	duplicate.ID = 0
	if err = s.Save(ctx, &duplicate); !errors.Is(err, ErrNewsConflict) {
		t.Fatalf("duplicate slug: %v", err)
	}
	changedSlug := a
	changedSlug.Slug = "changed"
	if err = s.Save(ctx, &changedSlug); !errors.Is(err, ErrNewsConflict) {
		t.Fatalf("changed permanent slug: %v", err)
	}
	stale := a
	a.Title = "Aktualizováno"
	if err = s.Save(ctx, &a); err != nil {
		t.Fatal(err)
	}
	if err = s.Save(ctx, &stale); !errors.Is(err, ErrNewsConflict) {
		t.Fatalf("stale edit: %v", err)
	}
	a, _ = s.ByID(ctx, a.ID)
	a.Status = "draft"
	if err = s.Save(ctx, &a); err != nil {
		t.Fatal(err)
	}
	if _, err = s.BySlug(ctx, a.Slug); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("withdrawn detail: %v", err)
	}
	if rows, err := s.Sitemap(ctx); err != nil || len(rows) != 0 {
		t.Fatalf("withdrawn sitemap: %v %v", rows, err)
	}
	a, _ = s.ByID(ctx, a.ID)
	a.Status = "published"
	if err = s.Save(ctx, &a); err != nil {
		t.Fatal(err)
	}
	a, _ = s.ByID(ctx, a.ID)
	if !a.PublishedAt.Equal(firstPublished) {
		t.Fatal("republication changed original timestamp")
	}
	// Future-dated rows must remain private even if introduced outside the editor.
	if _, err = db.Exec("UPDATE news_articles SET published_at=DATE_ADD(UTC_TIMESTAMP(), INTERVAL 1 DAY) WHERE id=?", a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err = s.BySlug(ctx, a.Slug); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("future detail: %v", err)
	}
	if rows, n, err := s.List(ctx, "", true, 12, 0); err != nil || n != 0 || len(rows) != 0 {
		t.Fatalf("future list: %v %v %v", rows, n, err)
	}
	if rows, err := s.Sitemap(ctx); err != nil || len(rows) != 0 {
		t.Fatalf("future sitemap: %v %v", rows, err)
	}
	if err = s.Delete(ctx, a.ID, a.Revision-1); !errors.Is(err, ErrNewsConflict) {
		t.Fatalf("stale delete: %v", err)
	}
	if err = s.Delete(ctx, a.ID, a.Revision); err != nil {
		t.Fatal(err)
	}
	if _, err = s.ByID(ctx, a.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("delete: %v", err)
	}
	// Deterministic ordering when publication times match; pagination and totals.
	for i := 0; i < 3; i++ {
		row := models.NewsArticle{Title: "Zpráva", Slug: fmt.Sprintf("zprava-%d", i), Summary: "P", Body: "B", Status: "published", CategoryID: cats[0].ID}
		if err = s.Save(ctx, &row); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = db.Exec("UPDATE news_articles SET published_at='2020-01-01 12:00:00'"); err != nil {
		t.Fatal(err)
	}
	page1, n, err := s.List(ctx, cats[0].Slug, true, 2, 0)
	if err != nil || n != 3 || len(page1) != 2 {
		t.Fatalf("page1: %v %v %v", page1, n, err)
	}
	page2, n, err := s.List(ctx, cats[0].Slug, true, 2, 2)
	if err != nil || n != 3 || len(page2) != 1 || page1[1].ID <= page2[0].ID {
		t.Fatalf("page2: %v %v %v", page2, n, err)
	}
	// Rollback only removes new tables, and migration can be applied again.
	if _, err = db.Exec("CREATE TABLE unrelated (id INT)"); err != nil {
		t.Fatal(err)
	}
	down, err := os.ReadFile("../../migrations/0012_news.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(down)); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO unrelated VALUES (1)"); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
}
