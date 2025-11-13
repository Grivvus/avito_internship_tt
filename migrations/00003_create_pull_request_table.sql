-- +goose Up
-- +goose StatementBegin

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    status pr_status NOT NULL DEFAULT 'OPEN',
    need_more_reviewers BOOLEAN NOT NULL DEFAULT FALSE,
    merged_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    
    -- Внешний ключ на автора
    CONSTRAINT fk_pull_requests_author 
        FOREIGN KEY (author_id) 
        REFERENCES users(id) 
        ON DELETE RESTRICT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE "pull_requests";

-- +goose StatementEnd
