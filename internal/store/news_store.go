package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/skybedy/laravel-tene.life/internal/models"
)

var ErrNewsConflict = errors.New("article changed or slug already exists")
var ErrNewsCategory = errors.New("category no longer exists")

type NewsStore struct{ DB *sql.DB }

// Public queries always share this predicate, including home, related and sitemap.
const newsPublished = "a.status = 'published' AND a.published_at IS NOT NULL AND a.published_at <= UTC_TIMESTAMP()"
const newsColumns = `a.id,a.category_id,a.title,a.slug,a.summary,a.body,a.image_url,a.image_alt,
a.source_name,a.source_url,a.seo_title,a.seo_description,a.status,a.published_at,a.created_at,a.updated_at,a.revision,c.id,c.name,c.slug`
const newsJoin = " FROM news_articles a JOIN news_categories c ON c.id=a.category_id "

func scanNews(row interface{ Scan(...any) error }) (a models.NewsArticle, err error) {
	err = row.Scan(&a.ID, &a.CategoryID, &a.Title, &a.Slug, &a.Summary, &a.Body, &a.ImageURL, &a.ImageAlt,
		&a.SourceName, &a.SourceURL, &a.SEOTitle, &a.SEODescription, &a.Status, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt, &a.Revision, &a.Category.ID, &a.Category.Name, &a.Category.Slug)
	return
}

func (s *NewsStore) Categories(ctx context.Context) ([]models.NewsCategory, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id,name,slug FROM news_categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.NewsCategory
	for rows.Next() {
		var c models.NewsCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (s *NewsStore) List(ctx context.Context, category string, public bool, limit, offset int) ([]models.NewsArticle, int, error) {
	where := " WHERE 1=1"
	var args []any
	if public {
		where += " AND " + newsPublished
	}
	if category != "" {
		where += " AND c.slug=?"
		args = append(args, category)
	}
	var count int
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*)"+newsJoin+where, args...).Scan(&count); err != nil {
		return nil, 0, err
	}
	order := " ORDER BY a.updated_at DESC,a.id DESC"
	if public {
		order = " ORDER BY a.published_at DESC,a.id DESC"
	}
	rows, err := s.DB.QueryContext(ctx, "SELECT "+newsColumns+newsJoin+where+order+" LIMIT ? OFFSET ?", append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []models.NewsArticle
	for rows.Next() {
		a, err := scanNews(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, a)
	}
	return result, count, rows.Err()
}

func (s *NewsStore) BySlug(ctx context.Context, slug string) (models.NewsArticle, error) {
	return scanNews(s.DB.QueryRowContext(ctx, "SELECT "+newsColumns+newsJoin+" WHERE a.slug=? AND "+newsPublished, slug))
}

func (s *NewsStore) ByID(ctx context.Context, id int64) (models.NewsArticle, error) {
	return scanNews(s.DB.QueryRowContext(ctx, "SELECT "+newsColumns+newsJoin+" WHERE a.id=?", id))
}

func newsWriteError(err error) error {
	var e *mysql.MySQLError
	if errors.As(err, &e) {
		if e.Number == 1062 {
			return ErrNewsConflict
		}
		if e.Number == 1452 {
			return ErrNewsCategory
		}
	}
	return err
}

func (s *NewsStore) Save(ctx context.Context, a *models.NewsArticle) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now().UTC().Truncate(time.Second)
	if a.ID != 0 {
		var slug string
		var published *time.Time
		var revision int
		err = tx.QueryRowContext(ctx, "SELECT slug,published_at,revision FROM news_articles WHERE id=? FOR UPDATE", a.ID).Scan(&slug, &published, &revision)
		if err != nil {
			return err
		}
		if a.Revision != revision || (published != nil && a.Slug != slug) {
			return ErrNewsConflict
		}
		a.PublishedAt = published
	} else {
		a.PublishedAt = nil
	}
	if a.Status == "published" && a.PublishedAt == nil {
		a.PublishedAt = &now
	}
	args := []any{a.CategoryID, a.Title, a.Slug, a.Summary, a.Body, a.ImageURL, a.ImageAlt, a.SourceName, a.SourceURL, a.SEOTitle, a.SEODescription, a.Status, a.PublishedAt, now}
	if a.ID == 0 {
		res, e := tx.ExecContext(ctx, `INSERT INTO news_articles (category_id,title,slug,summary,body,image_url,image_alt,source_name,source_url,seo_title,seo_description,status,published_at,updated_at,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, append(args, now)...)
		if e != nil {
			return newsWriteError(e)
		}
		a.ID, err = res.LastInsertId()
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE news_articles SET category_id=?,title=?,slug=?,summary=?,body=?,image_url=?,image_alt=?,source_name=?,source_url=?,seo_title=?,seo_description=?,status=?,published_at=?,updated_at=?,revision=revision+1 WHERE id=?`, append(args, a.ID)...)
	}
	if err != nil {
		return newsWriteError(err)
	}
	return tx.Commit()
}

func (s *NewsStore) Delete(ctx context.Context, id int64, revision int) error {
	res, err := s.DB.ExecContext(ctx, "DELETE FROM news_articles WHERE id=? AND revision=?", id, revision)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return ErrNewsConflict
	}
	return err
}

type NewsSitemapEntry struct {
	Slug      string
	UpdatedAt time.Time
}

func (s *NewsStore) Sitemap(ctx context.Context) ([]NewsSitemapEntry, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT a.slug,a.updated_at FROM news_articles a WHERE "+newsPublished+" ORDER BY a.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []NewsSitemapEntry
	for rows.Next() {
		var e NewsSitemapEntry
		if err := rows.Scan(&e.Slug, &e.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}
