-- +goose Up
-- +goose StatementBegin

CREATE TABLE teams (
    name VARCHAR(255) PRIMARY KEY
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE "teams";

-- +goose StatementEnd
