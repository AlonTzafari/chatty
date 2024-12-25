package main

import (
	"context"
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
		rows, err := db.QueryContext(c.Context(), `SELECT id, username FROM users WHERE username LIKE $1 LIMIT 10;`, usernameSearch+"%")
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

func isUserInChannel(userId string, channelId string, db *sql.DB, ctx context.Context) (bool, error) {
	row := db.QueryRowContext(ctx, `
			SELECT 1
			FROM channels_users
			WHERE channel_id = $1
			AND user_id = $2;`,
		channelId,
		userId,
	)
	var n uint8
	err := row.Scan(&n)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
