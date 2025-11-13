-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team VARCHAR(255),
    
    -- Внешний ключ на таблицу команд
    CONSTRAINT fk_users_team 
        FOREIGN KEY (team) 
        REFERENCES teams(name) 
        ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE "users";

-- +goose StatementEnd
