-- +goose Up
-- +goose StatementBegin

CREATE TABLE pull_request_reviewers (
    pr_id BIGINT NOT NULL,
    reviewer_id BIGINT NOT NULL,
    
    PRIMARY KEY (pr_id, reviewer_id)
    
    -- Уникальность: один пользователь не может быть назначен дважды на один PR
    CONSTRAINT uq_pull_request_reviewer 
        UNIQUE (pull_request_id, reviewer_id),
    
    CONSTRAINT fk_pr_reviewers_pull_request 
        FOREIGN KEY (pull_request_id) 
        REFERENCES pull_requests(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_pr_reviewers_reviewer 
        FOREIGN KEY (reviewer_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE "pull_request_reviewers"

-- +goose StatementEnd
