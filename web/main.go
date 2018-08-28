package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func checkHost(fn echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Host --->", c.Request().Host)
		fn(c)
		return nil
	}
}

func main() {
	// Echo instance
	e := echo.New()
	// Middleware

	e.Use(checkHost)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use()
	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
