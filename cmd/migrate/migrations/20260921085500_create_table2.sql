
-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS table2 
(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id int NOT NULL,
    description TEXT,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_services_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS table2;

-- +goose StatementEnd