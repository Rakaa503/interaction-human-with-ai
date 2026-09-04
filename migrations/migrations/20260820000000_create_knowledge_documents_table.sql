-- +goose Up

CREATE TABLE knowledge_documents (
    id BIGSERIAL PRIMARY KEY,

    title TEXT NOT NULL,

    content TEXT NOT NULL,

    url TEXT,

    source TEXT,

    category VARCHAR(100),

    content_hash CHAR(64) NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_knowledge_documents_source
    ON knowledge_documents(source);

CREATE INDEX idx_knowledge_documents_category
    ON knowledge_documents(category);

-- +goose Down

DROP TABLE knowledge_documents;