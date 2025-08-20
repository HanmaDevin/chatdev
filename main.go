package main

import (
	"time"

	"github.com/HanmaDevin/chatdev/types"
	"github.com/HanmaDevin/chatdev/views"
	"github.com/labstack/echo"
)

func indexHandler(c echo.Context) error {
	data := types.Index{
		AppName:     "My Echo App",
		CurrentTime: time.Now(),
	}
	component := views.Index(data)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func main() {
	e := echo.New()
	e.GET("/", indexHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
