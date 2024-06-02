package model

import "time"

type Title struct {
	Id              uint64    `db:"id"`
	Type            string    `db:"type"`
	Title           string    `db:"title"`
	TranslatedTitle string    `db:"translated_title"`
	PremierDate     time.Time `db:"premier_date"`
	Genre           string    `db:"genre"`
	AgeRating       string    `db:"age_rating"`
	Duration        string    `db:"duration"`
	Description     string    `db:"description"`
	Country         string    `db:"country"`

	NumberOfEpisodes int `db:"number_of_episodes"`
	NumberOfSeasons  int `db:"number_of_seasons"`
}

func GetAllTitles() ([]Title, error) {
	titles := []Title{}

	query := `SELECT
		t.id,
		tt.name AS type,
		t.title,
		t.translated_title,
		t.premier_date,
		t.genre,
		ar.text AS age_rating,
		t.duration,
		t.description,
		t.country,

		COALESCE(t.number_of_episodes, 0) AS number_of_episodes,
		COALESCE(t.number_of_seasons, 0) AS number_of_seasons
	FROM
		title AS t
	LEFT JOIN
		title_age_rating AS ar ON ar.id = t.age_rating
	LEFT JOIN
		title_type AS tt ON tt.type = t.type
	ORDER BY
		t.id;
	`

	if err := db.Select(&titles, query); err != nil {
		return nil, err
	}

	return titles, nil
}
