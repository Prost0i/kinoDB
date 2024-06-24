package model

import (
	"fmt"
	"strings"
	"time"
)

type Title struct {
	Id              uint64    `db:"id"`
	Type            string    `db:"type"`
	TypeChar        string    `db:"type_char"`
	Title           string    `db:"title"`
	TranslatedTitle string    `db:"translated_title"`
	PremierDate     time.Time `db:"premier_date"`
	Genre           string    `db:"genre"`
	AgeRatingId     uint64    `db:"age_rating_id"`
	AgeRating       string    `db:"age_rating"`
	Duration        string    `db:"duration"`
	Description     string    `db:"description"`
	Country         string    `db:"country"`

	NumberOfEpisodes int `db:"number_of_episodes"`
	NumberOfSeasons  int `db:"number_of_seasons"`

	RatingAvg     string `db:"rating_avg"`
	RatingCnt     int    `db:"rating_cnt"`
	RatingReviews int    `db:"rating_reviews"`

	DurationFormatted string
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

	t.DurationFormatted = strings.TrimSpace(durationStr)
}

func GetAllTitles() ([]Title, error) {
	titles := []Title{}

	query := `SELECT
		t.id,
		tt.name AS type,
		t.type AS type_char,
		t.title,
		t.translated_title,
		t.premier_date,
		t.genre,
		t.age_rating AS age_rating_id,
		ar.text AS age_rating,
		t.duration,
		t.description,
		t.country,

		t.number_of_episodes AS number_of_episodes,
		t.number_of_seasons AS number_of_seasons,

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
		t.type AS type_char,
		t.title,
		t.translated_title,
		t.premier_date,
		t.genre,
		t.age_rating AS age_rating_id,
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
		t.type AS type_char,
		t.title,
		t.translated_title,
		t.premier_date,
		t.genre,
		t.age_rating AS age_rating_id,
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

func UpdateTitle(titleId uint64, titleTitle, translatedTitle, typeChar, genre string,
	ageRatingId uint64,
	country, description, premierDate, duration string,
	numberOfEpisodes, numberOfSeasons int) (uint64, error) {
	query := `
		UPDATE title SET
			title = $1,
			translated_title = $2,
			type = $3,
			genre = $4,
			age_rating = $5,
			country = $6,
			description = $7,
			premier_date = $8,
			duration = $9,
			number_of_episodes = $10,
			number_of_seasons = $11
		WHERE
			title.id = $12
		RETURNING id
	`
	var updatedTitleId uint64
	err := db.Get(
		&updatedTitleId,
		query,
		titleTitle,
		translatedTitle,
		typeChar,
		genre,
		ageRatingId,
		country,
		description,
		premierDate,
		duration,
		numberOfEpisodes,
		numberOfSeasons,
		titleId)
	if err != nil {
		return 0, err
	}

	return updatedTitleId, nil
}

func InsertTitle(titleTitle, translatedTitle, typeChar, genre string, ageRatingId uint64,
	country, description, premierDate, duration string,
	numberOfEpisodes, numberOfSeasons int) (uint64, error) {

	query := `
		INSERT INTO title(
			title,
			translated_title,
			type,
			genre,
			age_rating,
			country,
			description,
			premier_date,
			duration,
			number_of_episodes,
			number_of_seasons)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	var titleId uint64
	err := db.Get(
		&titleId,
		query,
		titleTitle,
		translatedTitle,
		typeChar,
		genre,
		ageRatingId,
		country,
		description,
		premierDate,
		duration,
		numberOfEpisodes,
		numberOfSeasons)
	if err != nil {
		return 0, err
	}

	return titleId, nil
}

func DeleteTitle(titleId uint64) error {
	query := `
		DELETE FROM title WHERE title.id = $1
	`

	_, err := db.Exec(query, titleId)
	if err != nil {
		return err
	}

	return nil
}
