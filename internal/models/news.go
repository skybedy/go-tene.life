package models

import "time"

type NewsCategory struct {
	ID         int64
	Name, Slug string
}

type NewsArticle struct {
	ID, CategoryID                            int64
	Title, Slug, Summary, Body                string
	ImageURL, ImageAlt, SourceName, SourceURL string
	SEOTitle, SEODescription, Status          string
	PublishedAt                               *time.Time
	CreatedAt, UpdatedAt                      time.Time
	Revision                                  int
	Category                                  NewsCategory
}

func (a NewsArticle) PublishedDate() string {
	if a.PublishedAt == nil {
		return ""
	}
	loc, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		loc = time.UTC
	}
	return a.PublishedAt.In(loc).Format("2. 1. 2006 15:04")
}

func (a NewsArticle) PublishedISO() string {
	if a.PublishedAt == nil {
		return ""
	}
	return a.PublishedAt.UTC().Format(time.RFC3339)
}
