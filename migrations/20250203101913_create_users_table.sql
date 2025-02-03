-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       first_name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100) NOT NULL,
                       balance INT NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                              type VARCHAR(50) NOT NULL,
                              amount INT NOT NULL,
                              from_id INT REFERENCES users(id) ON DELETE SET NULL,
                              to_id INT REFERENCES users(id) ON DELETE SET NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
