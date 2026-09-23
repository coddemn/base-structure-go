
-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    login VARCHAR(255) NOT NULL,
    pass VARCHAR(255) NOT NULL
)

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd