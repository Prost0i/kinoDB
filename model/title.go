package model

import (
	"fmt"
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
		review_rating AS r ON r.title_id = t.id
	GROUP BY
		t.id,
		tt.name,
		ar.text
	ORDER BY
		t.id DESC;
	`

	if err := db.Select(&titles, query); err != nil {
		return nil, err
	}

	return titles, nil
}

func FilterTitles(title, genre, typeChar, orderBy string) ([]Title, error) {
	titles := []Title{}

	b := strings.Builder{}
	b.WriteString(
		`SELECT
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
		review_rating AS r ON r.title_id = t.id
	`)

	if title != "" || genre != "" || typeChar != "all" {
		b.WriteString(`
			WHERE
		`)
	}

	if typeChar != "all" {
		b.WriteString(fmt.Sprintf(`
				(t.type = '%s')
				`, typeChar))
	}

	if typeChar != "all" && title != "" {
		b.WriteString(`
		AND
		`)
	}

	if title != "" {
		b.WriteString(fmt.Sprintf(`
				(LOWER(t.title) LIKE LOWER('%%%s%%')
				OR LOWER(t.translated_title) LIKE LOWER('%%%s%%'))
		`, title, title))
	}

	if (title != "" && genre != "") ||
		(typeChar != "all" && genre != "") {
		b.WriteString(`
		AND
		`)
	}

	if genre != "" {
		genreSplit := strings.Split(genre, ",")
		for i := range genreSplit {
			genreSplit[i] = strings.ReplaceAll(genreSplit[i], " ", "")
		}

		b.WriteString(`
				(to_tsvector(LOWER(t.genre)) @@ to_tsquery('`)

		for i := range genreSplit {
			b.WriteString(genreSplit[i])
			if i < len(genreSplit)-1 {
				b.WriteString(" & ")
			}
		}

		b.WriteString(`'))`)
	}

	b.WriteString(`
	GROUP BY
		t.id,
		tt.name,
		ar.text`)

	if orderBy == "rating" {
		b.WriteString(fmt.Sprintf(`
			ORDER BY
			rating_avg DESC;
		`))
	} else if orderBy == "rating_cnt" {
		b.WriteString(fmt.Sprintf(`
			ORDER BY
			rating_cnt DESC;
		`))
	} else if orderBy == "title" {
		b.WriteString(fmt.Sprintf(`
			ORDER BY
			t.title;
		`))
	} else {
		b.WriteString(fmt.Sprintf(`
			ORDER BY
			t.id;
		`))
	}

	query := b.String()
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
		review_rating AS r ON r.title_id = t.id
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
