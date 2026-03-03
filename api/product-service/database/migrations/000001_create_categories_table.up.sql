CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    parent_id BIGINT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(255) NOT NULL,
    status BOOLEAN DEFAULT TRUE,
    slug VARCHAR(120) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
CREATE INDEX idx_categories_status ON categories(status);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE UNIQUE INDEX categories_slug_deletedat_key ON categories (slug)
WHERE deleted_at IS NULL;