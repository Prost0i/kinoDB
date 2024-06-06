-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE review_rating(
	id BIGSERIAL PRIMARY KEY,
	value INT NOT NULL CHECK (value >= 1 AND value <= 10),
	title_id SERIAL NOT NULL REFERENCES title ON DELETE CASCADE,
	site_user_id SERIAL NOT NULL REFERENCES site_user ON DELETE CASCADE,
	review_title VARCHAR(500) DEFAULT NULL,
	review TEXT DEFAULT NULL
);

CREATE TABLE title_review_rating(
	title_id SERIAL NOT NULL REFERENCES title ON DELETE CASCADE,
	review_rating_id SERIAL NOT NULL UNIQUE REFERENCES review_rating ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS title_review_rating;
DROP TABLE IF EXISTS review_rating;
DELETE FROM site_user WHERE id = 228;
