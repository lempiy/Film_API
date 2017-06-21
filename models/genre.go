package models

import (
	"database/sql"
	"fmt"
	"github.com/lempiy/echo_api/types"
)

type genre struct{}
var Genre *genre

func (g *genre) ReadByFilmID(filmID int) ([]types.Genre, error) {
	var genres []types.Genre
	var genre types.Genre
	var rows *sql.Rows

	querySQL := `SELECT g.id, g.name, g.added_at
                FROM film_genre f
                  LEFT JOIN genre g ON g.id=f.genre_id
                WHERE f.film_id=$1 ;`

	rows = Database.Query(querySQL, filmID)

	defer rows.Close()
	for rows.Next() {
		genre = types.Genre{}

		err = rows.Scan(&genre.ID, &genre.Name, &genre.AddedAt)
		if err != nil {
			fmt.Println(err)
		}

		genres = append(genres, genre)
	}

	return genres, nil
}
