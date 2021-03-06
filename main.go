package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lempiy/echo_api/views/user"
	"net/http"
	"github.com/lempiy/echo_api/views/film"
	"os"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8001"
	}
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Film API 2017.")
	})
	e.POST("/login", user.Login)
	e.POST("/auth", user.Register)
	e.POST("/api/v1/film", film.Add, middleware.JWT([]byte("secret")))
	e.POST("/api/v1/film/rent", film.Rent, middleware.JWT([]byte("secret")))
	e.POST("/api/v1/film/finish", film.FinishRent, middleware.JWT([]byte("secret")))
	e.GET("/api/v1/film", film.Get)
	e.GET("/api/v1/rented-film", film.GetRentedFilms, middleware.JWT([]byte("secret")))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	//Test group

	r := e.Group("/test")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", user.Test)
	e.Logger.Fatal(e.Start(":"+PORT))
}
