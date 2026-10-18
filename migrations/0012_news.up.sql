CREATE TABLE news_categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY news_category_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO news_categories (name, slug) VALUES
('Doprava', 'doprava'),
('Počasí a příroda', 'pocasi-a-priroda'),
('Život na Tenerife', 'zivot-na-tenerife'),
('Akce a kultura', 'akce-a-kultura'),
('Turistika', 'turistika');

CREATE TABLE news_articles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  category_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(200) NOT NULL,
  slug VARCHAR(180) NOT NULL,
  summary VARCHAR(500) NOT NULL,
  body MEDIUMTEXT NOT NULL,
  image_url VARCHAR(2048) NOT NULL DEFAULT '',
  image_alt VARCHAR(200) NOT NULL DEFAULT '',
  source_name VARCHAR(200) NOT NULL DEFAULT '',
  source_url VARCHAR(2048) NOT NULL DEFAULT '',
  seo_title VARCHAR(200) NOT NULL DEFAULT '',
  seo_description VARCHAR(500) NOT NULL DEFAULT '',
  status ENUM('draft', 'published') NOT NULL DEFAULT 'draft',
  published_at DATETIME NULL COMMENT 'First publication, UTC; retained when withdrawn',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  revision INT NOT NULL DEFAULT 1,
  PRIMARY KEY (id),
  UNIQUE KEY news_article_slug (slug),
  KEY news_publication (status, published_at, id),
  KEY news_category_publication (category_id, status, published_at, id),
  CONSTRAINT news_article_category FOREIGN KEY (category_id) REFERENCES news_categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
