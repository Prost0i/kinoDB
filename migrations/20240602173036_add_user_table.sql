-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE site_user(
	id BIGSERIAL PRIMARY KEY,
	email VARCHAR(500) NOT NULL UNIQUE,
	username VARCHAR(500) NOT NULL,
	password_hash VARCHAR(500) NOT NULL,
	is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS site_user;
