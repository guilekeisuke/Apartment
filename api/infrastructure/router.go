package infrastructure

import (
	"api/handler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.POST("/pre_signup", handler.PreSignup)
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)
	e.GET("/movie_upcoming", handler.MovieUpcoming)
	// e.GET("/movie_list", handler.getTmdbInfo)

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(handler.Config))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
