package model

import (
	"strings"
	"time"
)

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

	RatingAvg     string `db:"rating_avg"`
	RatingCnt     int    `db:"rating_cnt"`
	RatingReviews int    `db:"rating_reviews"`
}

func (t *Title) ConvertDuration() {
	durationSplit := strings.Split(t.Duration, ":")
	durationStr := ""
	timeLabels := []string{" ч. ", " м. ", " с. "}
	for i := range durationSplit {
		if durationSplit[i] != "00" {
			durationStr += strings.TrimLeft(durationSplit[i], "0") + timeLabels[i]
		}
	}

	t.Duration = strings.TrimSpace(durationStr)
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
		COALESCE(t.number_of_seasons, 0) AS number_of_seasons,

		COALESCE(AVG(r.value), 0.0)::numeric(10, 1) as rating_avg,
		COALESCE(COUNT(r), 0) as rating_cnt,
		COALESCE(COUNT(r) FILTER (WHERE review IS NOT NULL), 0) AS rating_reviews
	FROM
		title AS t
	LEFT JOIN
		title_age_rating AS ar ON ar.id = t.age_rating
	LEFT JOIN
		title_type AS tt ON tt.type = t.type
	LEFT OUTER JOIN
		title_rating AS tr ON t.id = tr.title_id
	LEFT OUTER JOIN
		rating AS r ON r.id = tr.rating_id
	GROUP BY
		t.id,
		tt.name,
		ar.text
	ORDER BY
		t.id;
	`

	if err := db.Select(&titles, query); err != nil {
		return nil, err
	}

	return titles, nil
}

func GetTitleById(id uint64) (Title, error) {
	title := Title{}

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
		COALESCE(t.number_of_seasons, 0) AS number_of_seasons,

		COALESCE(AVG(r.value), 0)::numeric(10, 1) AS rating_avg,
		COALESCE(COUNT(r), 0) AS rating_cnt,
		COALESCE(COUNT(r) FILTER (WHERE review IS NOT NULL), 0) AS rating_reviews
	FROM
		title AS t
	LEFT JOIN
		title_age_rating AS ar ON ar.id = t.age_rating
	LEFT JOIN
		title_type AS tt ON tt.type = t.type
	LEFT OUTER JOIN
		title_rating AS tr ON t.id = tr.title_id
	LEFT OUTER JOIN
		rating AS r ON r.id = tr.rating_id
	WHERE
		t.id = $1
	GROUP BY
		t.id,
		tt.name,
		ar.text
	ORDER BY
		t.id;
	`

	if err := db.Get(&title, query, id); err != nil {
		return Title{}, err
	}

	return title, nil
}
