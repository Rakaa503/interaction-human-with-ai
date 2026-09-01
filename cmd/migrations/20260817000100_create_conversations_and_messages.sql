-- +goose Up

CREATE TABLE conversations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_conversations_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_conversations_user_id
    ON conversations(user_id);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_messages_conversation
        FOREIGN KEY (conversation_id)
        REFERENCES conversations(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_messages_role
        CHECK (role IN ('system', 'user', 'assistant'))
);

CREATE INDEX idx_messages_conversation_id
    ON messages(conversation_id);

CREATE INDEX idx_messages_created_at
    ON messages(created_at);

-- +goose Down

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;