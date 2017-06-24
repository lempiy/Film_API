package models

import (
	"github.com/lempiy/echo_api/types"
	"fmt"
	"strconv"
)

type film struct{}
var Film *film

func (f *film) Create(film *types.Film, genresIDs []int) error {
	sqlQuery := `INSERT INTO film(name,year,added_at)
			VALUES($1, $2, now())`
	ID := Database.InsertWithReturningID(sqlQuery, film.Name, film.Year)
	return f.postFilmGenres(ID, genresIDs)
}

func (f *film) Read(params *types.GetFilmParams) ([]types.Film, bool, int, error) {
	yearFilter := f.yearFilter(params.Year, 2, true)
	var sqlQuery string
	if len(params.Genre) != 0 {
		sqlQuery = f.queryGenreFilter(yearFilter, params.Limit, params.Genre)
	} else {
		sqlQuery =
			fmt.Sprintf(`WITH f AS (SELECT *, COUNT(*) OVER () AS total_items
         		FROM film f
         		%s)
				SELECT f.id, f.name, f.year, f.added_at, total_items
					FROM f
				ORDER BY f.added_at ASC
			LIMIT %s OFFSET $1;`, yearFilter, f.getLimit(params.Limit))
	}
	return f.getFilms(params, sqlQuery)
}

//ReadByID returns pointer to types.Film by its ID
func (f *film) ReadByID(filmID int) (*types.Film, error) {
	var film types.Film
	sqlQuery := `SELECT f.id, f.name, f.year, f.added_at
				FROM film f WHERE id=$1;`
	rows := Database.Query(sqlQuery, filmID)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&film.ID,&film.Name,&film.Year,&film.AddedAt)
		if err != nil {
			return nil, err
		}
	}
	return &film, nil
}

type rent struct {
	ID int
	FilmID int
	UserID int
}

func (f *film) getFilms(params *types.GetFilmParams, query string, arguments ...interface{}) ([]types.Film, bool, int, error) {
	films := make([]types.Film, 0)
	var film types.Film
	var count int
	agrs := []interface{}{
		params.Offset,
	}
	agrs = append(agrs, arguments...)

	if params.Year != 0 {
		agrs = append(agrs, params.Year)
	}
	fmt.Print(query)
	rows := Database.Query(query, agrs...)
	defer rows.Close()
	for rows.Next() {
		film = types.Film{}
		err = rows.Scan(&film.ID, &film.Name, &film.Year, &film.AddedAt, &count)
		if err != nil {
			return nil, false, 0, err
		}
		genres, err := Genre.ReadByFilmID(film.ID)
		if err != nil {
			fmt.Printf("Error while getting genres for film - %s\n", film.ID)
			return nil, false, 0, err
		}
		film.Genres = genres
		films = append(films, film)
	}
	left := count - (params.Offset + len(films)) > 0
	return films, left, count, nil
}

func (f *film) ReadRentedFilms(userID int, params *types.GetFilmParams) ([]types.Film, bool, int, error) {
	yearFilter := f.yearFilter(params.Year, 3, false)
	var sqlQuery string
	if params.Genre != "" {
		sqlQuery = f.queryGenreFilterForRented(yearFilter, params.Limit, params.Genre)
	} else {
		sqlQuery =
			fmt.Sprintf(`WITH r AS (SELECT *, COUNT(*) OVER () AS total_items
					FROM rented_film r
						LEFT JOIN film f ON r.film_id=f.id
					WHERE r.user_id=$2 %s)
					SELECT f.id, f.name, f.year, f.added_at, total_items
						FROM r
						LEFT JOIN film f ON r.film_id=f.id
					ORDER BY f.added_at ASC
					LIMIT %s OFFSET $1;`, yearFilter, f.getLimit(params.Limit))
	}
	return f.getFilms(params, sqlQuery, userID)
}

func (f *film) queryGenreFilter(yearFilter string, limit int, genres string) string {
	sqlQuery :=
		fmt.Sprintf(`WITH f AS (SELECT DISTINCT(g.film_id), f.id, f.name, f.year, f.added_at,
					COUNT(*) OVER () AS total_items
				FROM film_genre g
					LEFT JOIN film f ON g.film_id=f.id
				WHERE g.genre_id IN (%s)
					%s)
				SELECT f.id, f.name, f.year, f.added_at, total_items
						FROM f
				ORDER BY f.added_at ASC
				LIMIT %s OFFSET $1;`, genres, yearFilter, f.getLimit(limit))
	return sqlQuery
}

func (f *film) queryGenreFilterForRented(yearFilter string, limit int, genres string) string {
	sqlQuery :=
		fmt.Sprintf(`WITH f AS (SELECT DISTINCT(g.film_id), f.id, f.name, f.year, f.added_at,
					COUNT(*) OVER () AS total_items
				FROM film_genre g
					LEFT JOIN film f ON g.film_id=f.id
					LEFT JOIN rented_film r ON r.film_id=g.film_id
				WHERE g.genre_id IN (%s) AND r.user_id=$2
					%s)
				SELECT f.id, f.name, f.year, f.added_at, total_items
						FROM f
				ORDER BY f.added_at ASC
				LIMIT %s OFFSET $1;`, genres, yearFilter, f.getLimit(limit))
	return sqlQuery
}

func (f *film) getLimit(limit int) string {
	if limit > 0 {
		return strconv.Itoa(limit)
	}
	return "ALL"
}

func (f *film) yearFilter(year int, paramOrder int, initial bool) string {
	if year == 0 {
		return ""
	}
	if initial {
		return fmt.Sprintf(`WHERE f.year=$%d`, paramOrder)
	}
	return fmt.Sprintf(`AND f.year=$%d`, paramOrder)
}

//ReadRentByID returns pointer to rent struct by film ID
func (f *film) readRentByFilmAndUser(filmID int, userID int) (*rent, error) {
	var film rent
	sqlQuery := `SELECT r.id, r.film_id, r.user_id
				FROM rented_film r WHERE film_id=$1 AND user_id=$2;`
	rows := Database.Query(sqlQuery, filmID, userID)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&film.ID,&film.FilmID,&film.UserID)
		if err != nil {
			return nil, err
		}
	}
	return &film, nil
}

//postFilmGenres creates genres instances for some new film
func (f *film) postFilmGenres(filmID int, genresIDs []int) error {
	sqlQuery := `INSERT INTO film_genre(film_id, genre_id, added_at) VALUES`
	var err error
	for i, id := range genresIDs {
		var queryValue string
		if i == 0 {
			queryValue = fmt.Sprintf("(%d, %d, now())", filmID, id)
		} else {
			queryValue = fmt.Sprintf(", (%d, %d, now())", filmID, id)
		}
		sqlQuery += queryValue
		if err != nil {
			return err
		}
	}
	sqlQuery += `;`
	return Database.SingleQuery(sqlQuery)
}

//Rent creates a rented film data, if film doesn't exist's in DB it returns true as a first argument
func(f *film) Rent(filmID int, userID int) (notExist bool, err error) {
	sqlQuery := `INSERT INTO rented_film(film_id, user_id, added_at) VALUES($1, $2, now())`
	film, err := f.ReadByID(filmID)
	if err != nil {
		notExist = false
		return
	}
	if film.ID == 0 {
		notExist = true
		return
	}
	rent, err := f.readRentByFilmAndUser(filmID, userID)
	if err != nil {
		notExist = false
		return
	}
	if rent.ID == 0 {
		notExist = true
		return
	}
	err = Database.SingleQuery(sqlQuery, filmID, userID)
	notExist = false
	return
}

//FinishRent sets renting object to finished, if film was never
// rented by user before it returns true as a first argument
func(f *film) FinishRent(filmID int, userID int) (notExist bool, err error) {
	sqlQuery := `UPDATE rented_film
				SET finished = 1
				WHERE film_id=$1 AND user_id=$2`

	rent, err := f.readRentByFilmAndUser(filmID, userID)

	if err != nil {
		return
	}

	if rent.ID == 0 {
		notExist = true
		return
	}
	err = Database.SingleQuery(sqlQuery, filmID, userID)
	return
}

func sliceItoa(sliceOfInt []int) (sliceOfStr []string) {
	for _, n := range sliceOfInt {
		sliceOfStr = append(sliceOfStr, strconv.Itoa(n))
	}
	return
}