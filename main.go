package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/HanmaDevin/chatdev/db"
	"github.com/HanmaDevin/chatdev/types"
	"github.com/HanmaDevin/chatdev/views"
	"github.com/labstack/echo"
)

var currentTime = time.Now().Format("15:04") + " Uhr"

func indexHandler(c echo.Context) error {
	data := types.Index{
		AppName:     "ChatDev",
		CurrentTime: currentTime,
	}
	component := views.Index(data)
	fmt.Printf("Data: %+v\n", data)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func loginHandler(c echo.Context) error {
	formData := types.FormData{}
	component := views.Form(formData)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func loginAuthHandler(c echo.Context) error {
	formData := types.FormData{}
	formData.Username = c.FormValue("username")
	formData.Password = c.FormValue("password")

	_, err := db.ValidUser(formData.Username, formData.Password)
	if err != nil {
		formData.Errors.UsernameError = "Invalid username or password"
		component := views.Form(formData)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	return c.Redirect(301, "/chats")
}

func chatsHandler(c echo.Context) error {

	data := db.GetAllChats()
	component := views.Chats(data)
	fmt.Printf("Data: %+v\n", data)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func convHandler(c echo.Context) error {
	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	fmt.Printf("ID: %s\n, ID Int: %d\n, Error: %v\n", id, idint, err)
	// database simulation
	chat, err := db.GetChatByID(idint)
	component := views.Conversation(chat)

	return component.Render(c.Request().Context(), c.Response().Writer)
}

func main() {
	e := echo.New()

	e.GET("/", indexHandler)
	e.GET("/chats", chatsHandler)
	e.GET("/login", loginHandler)
	e.POST("/login", loginAuthHandler)
	e.GET("/chat/:id", convHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
