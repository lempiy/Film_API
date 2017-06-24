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
	var films []types.Film
	var film types.Film
	var count int
	limit := "ALL"
	if params.Limit != 0 {
		limit = strconv.Itoa(params.Limit)
	}
	sqlQuery := fmt.Sprintf(
		`WITH f AS (SELECT *, COUNT(*) OVER () AS total_items
         		FROM film)
				SELECT f.id, f.name, f.year, f.added_at, total_items
					FROM f
				ORDER BY f.added_at ASC
				LIMIT %s OFFSET %d;`, limit, params.Offset)
	rows := Database.Query(sqlQuery)
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
	fmt.Println("count:", count)
	left := count - (params.Offset + len(films)) > 0
	return films, left, count, nil
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

func (f *film) ReadRentedFilms(userID int, params *types.GetFilmParams) ([]types.Film, bool, int, error) {
	var films []types.Film
	var film types.Film
	var count int
	limit := "ALL"
	if params.Limit != 0 {
		limit = strconv.Itoa(params.Limit)
	}
	sqlQuery := fmt.Sprintf(
		`WITH r AS (SELECT *, COUNT(*) OVER () AS total_items
         		FROM rented_film
         		WHERE user_id=$1)
				SELECT f.id, f.name, f.year, f.added_at, total_items
					FROM r
					LEFT JOIN film f ON r.film_id=f.id
				ORDER BY f.added_at ASC
				LIMIT %s OFFSET %d;`, limit, params.Offset)
	rows := Database.Query(sqlQuery, userID)
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
	fmt.Println("count:", count)
	left := count - (params.Offset + len(films)) > 0
	return films, left, count, nil
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