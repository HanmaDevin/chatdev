package main

import (
	"context"

	"github.com/HanmaDevin/chatdev/dto"
	"github.com/HanmaDevin/chatdev/templates"
	"github.com/HanmaDevin/chatdev/types"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

var (
	password, _ = bcrypt.GenerateFromPassword([]byte("admin"), 8)
	user        = dto.UserDTO{
		ID:       uuid.New().String(),
		Username: "admin",
		Password: string(password),
	}
)

func main() {
	e := echo.New()
	manager := types.NewManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go manager.HandleClientListEventChan(ctx)

	e.GET("/", func(c echo.Context) error {
		component := templates.Index()
		return component.Render(ctx, c.Response().Writer)
	})

	authChat := e.Group("/chat")
	authChat.GET("", func(c echo.Context) error {
		component := templates.Chat()
		return component.Render(ctx, c.Response().Writer)
	})

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if username == "" || password == "" {
			return c.String(400, "Username and password are required")
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil && user.Username == username {
			c.Set("user", user)
			return c.Redirect(301, "/chat")
		} else {
			return c.String(401, "Invalid username or password")
		}
	})

	e.GET("/chatroom", func(c echo.Context) error {
		return manager.Handle(c, ctx)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
