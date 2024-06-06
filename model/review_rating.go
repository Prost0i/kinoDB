package model

import (
	"database/sql"
)

type ReviewRating struct {
	Id          uint64 `db:"id"`
	Value       int    `db:"value"`
	SiteUserId  uint64 `db:"site_user_id"`
	TitleId     uint64 `db:"title_id"`
	ReviewTitle string `db:"review_title"`
	Review      string `db:"review"`

	Username string `db:"username"`
}

func GetAllReviewsForTitleByTitleId(titleId uint64) ([]ReviewRating, error) {
	var reviews []ReviewRating
	query := `
		SELECT
			review_rating.id,
			review_rating.value,
			review_rating.site_user_id,
			review_rating.title_id,
			review_rating.review_title,
			review_rating.review,
			site_user.username
		FROM
			review_rating
		JOIN
			site_user ON site_user.id = review_rating.site_user_id
		WHERE
			review_rating.title_id = $1
			AND
			(review_rating.review_title, review_rating.review) IS NOT NULL
	`

	if err := db.Select(&reviews, query, titleId); err != nil {
		return nil, err
	}

	return reviews, nil
}

func GetReviewRatingByUserId(titleId, userId uint64) (bool, ReviewRating, error) {
	reviewRating := ReviewRating{}
	query := `
		SELECT
			review_rating.id,
			review_rating.value,
			review_rating.site_user_id,
			review_rating.title_id,
			COALESCE(review_rating.review_title, '') as review_title,
			COALESCE(review_rating.review, '') as review,
			site_user.username
		FROM
			review_rating
		JOIN
			site_user ON site_user.id = review_rating.site_user_id
		WHERE
			review_rating.title_id = $1
			AND
			review_rating.site_user_id = $2
	`

	if err := db.Get(&reviewRating, query, titleId, userId); err != nil {
		if err == sql.ErrNoRows {
			return false, ReviewRating{}, nil
		}

		return false, ReviewRating{}, err
	}

	return true, reviewRating, nil

}

func InsertOnlyRating(rating int, titleId, userId uint64) (uint64, error) {
	query := `
		INSERT INTO
			review_rating(id, value, title_id, site_user_id)
		VALUES
			(DEFAULT, $1, $2, $3) RETURNING id;
	`

	var reviewRatingId uint64
	if err := db.Get(&reviewRatingId, query, rating, titleId, userId); err != nil {
		return 0, err
	}

	return reviewRatingId, nil
}

func InsertReview(rating int, reviewTitle, review string, titleId, userId uint64) (uint64, error) {
	query := `
		INSERT INTO
			review_rating(id, value, title_id, site_user_id, review_title, review)
		VALUES
			(DEFAULT, $1, $2, $3, $4, $5) RETURNING id;
	`

	var reviewRatingId uint64
	if err := db.Get(&reviewRatingId, query, rating, titleId, userId, reviewTitle, review); err != nil {
		return 0, err
	}

	return reviewRatingId, nil
}

func UpdateReviewRating(rating int, reviewTitle, review string, titleId, userId uint64) (uint64, error) {
	var query string
	var reviewRatingId uint64
	if reviewTitle == "" && review == "" {
		query = `
		UPDATE review_rating
			SET value = $1, review_title = NULL, review = NULL
		WHERE
			title_id = $2 AND site_user_id = $3
		RETURNING id
		`

		if err := db.Get(&reviewRatingId, query, rating, titleId, userId); err != nil {
			return 0, err
		}

		return reviewRatingId, nil
	} else {
		query = `
		UPDATE review_rating
			SET value = $1, review_title = $2, review = $3
		WHERE
			title_id = $4 AND site_user_id = $5
		RETURNING id
		`

		if err := db.Get(&reviewRatingId, query, rating, reviewTitle, review, titleId, userId); err != nil {
			return 0, err
		}

		return reviewRatingId, nil
	}
}

func DeleteReviewRating(titleId, userId uint64) error {
	query := `
		DELETE FROM review_rating WHERE title_id = $1 AND site_user_id = $2
	`

	_, err := db.Exec(query, titleId, userId)
	if err != nil {
		return err
	}

	return nil
}
