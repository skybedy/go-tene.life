-- Explicit rollback deletes the news section and its content only.
DROP TABLE IF EXISTS news_articles;
DROP TABLE IF EXISTS news_categories;
