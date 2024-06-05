-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE rating(
	id BIGSERIAL PRIMARY KEY,
	value INT NOT NULL CHECK (value >= 1 AND value <= 10),
	site_user_id SERIAL NOT NULL REFERENCES site_user,
	review TEXT
);

CREATE TABLE title_rating(
	title_id SERIAL NOT NULL REFERENCES title,
	rating_id SERIAL NOT NULL UNIQUE REFERENCES rating
);

INSERT INTO site_user
	VALUES (228, 'example', 'example@example.com', 'no_password_no_hash');

INSERT INTO rating 
	VALUES
	(1, 7, 228, NULL),
	(2, 10, 228, 'This is my review'),
	(3, 7, 228, 'Another review');

INSERT INTO title_rating
	VALUES
	(1, 1),
	(1, 2),
	(2, 3);



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS title_rating;
DROP TABLE IF EXISTS rating;
DELETE FROM site_user WHERE id = 228;
