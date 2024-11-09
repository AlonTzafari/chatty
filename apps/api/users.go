package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Users(app *fiber.App, db *sql.DB) {
	app.Get("/api/users", func(c *fiber.Ctx) error {
		authCtx, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.SendStatus(http.StatusUnauthorized)
		}
		usernameSearch := c.Query("username")
		rows, err := db.Query(`SELECT id, username FROM users WHERE username LIKE $1 LIMIT 10;`, usernameSearch+"%")
		if err != nil {
			log.Print("ERROR Query", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		var users []User
		for rows.Next() {
			var (
				id       string
				username string
			)
			err = rows.Scan(&id, &username)
			if err != nil {
				log.Print("ERROR Scan", err)
				return c.SendStatus(http.StatusInternalServerError)
			}
			if id == authCtx.UserId {
				continue
			}
			users = append(users, User{id, username})
		}
		return c.JSON(users)
	})
}
