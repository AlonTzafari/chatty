package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Channel struct {
	Id        string
	Name      string
	Avatar    string
	CreatedAt string
}

type CreateChannelReq struct {
	Name    string   `json:"name" xml:"name" form:"name"`
	Members []string `json:"members" xml:"members" form:"members"`
}

func Channels(app *fiber.App, db *sql.DB) {
	app.Post("/api/channels", func(c *fiber.Ctx) error {
		_, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.SendStatus(http.StatusUnauthorized)
		}
		var createChannelReq CreateChannelReq
		err := c.BodyParser(&createChannelReq)
		if err != nil || createChannelReq.Name == "" || len(createChannelReq.Members) == 0 {
			log.Print(err)
			return c.SendStatus(http.StatusBadRequest)
		}
		tx, err := db.Begin()
		if err != nil {
			log.Print(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		defer tx.Rollback()
		row := tx.QueryRow(`INSERT INTO channels (name, avatar, createdAt) VALUES ($1, $2, $3) RETURNING id;`, createChannelReq.Name, nil, time.Now().UTC().Format(time.RFC3339))
		var channelId string
		err = row.Scan(&channelId)
		if err != nil {
			log.Print("CHANNEL INSERT", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		sqlStr := `INSERT INTO channels_users (channel_id, user_id) VALUES `
		var vals []any
		for i, userId := range createChannelReq.Members {
			rowStr := fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
			if i < len(createChannelReq.Members)-1 {
				rowStr += ","
			}
			sqlStr += rowStr
			vals = append(vals, channelId, userId)
		}
		sqlStr += ";"
		_, err = tx.Exec(sqlStr, vals...)
		if err != nil {
			log.Print("CHANNEL_USER INSERT ", sqlStr, err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		err = tx.Commit()
		if err != nil {
			log.Print(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.Status(http.StatusCreated).JSON(struct {
			Id string `json:"id"`
		}{channelId})

	})
	app.Get("/api/channels", func(c *fiber.Ctx) error {
		userId := c.Query("user_id")
		if userId == "" {
			return c.SendStatus(http.StatusBadRequest)
		}
		auth, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.SendStatus(http.StatusUnauthorized)
		}
		if auth.UserId != userId {
			return c.SendStatus(http.StatusForbidden)
		}
		rows, err := db.Query(`
			SELECT c.id, c.name, c.avatar, c.createdAt
			FROM channels c
			JOIN channels_users cu ON c.id = cu.channel_id
			JOIN users u ON u.id = cu.user_id
			WHERE u.id = $1;
		`, userId)
		if err != nil {
			log.Println(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		var channels []Channel
		for rows.Next() {
			var (
				id        string
				name      string
				avatar    *string
				createdAt string
			)
			err := rows.Scan(&id, &name, &avatar, &createdAt)
			if err != nil {
				log.Println(err)
				return c.SendStatus(http.StatusInternalServerError)
			}
			emptyStr := ""
			if avatar == nil {
				avatar = &emptyStr
			}
			channels = append(channels, Channel{id, name, *avatar, createdAt})
		}
		return c.JSON(channels)
	})
	app.Get("/api/channels/:id", func(c *fiber.Ctx) error {
		channelId := c.Params("id")
		if channelId == "" {
			return c.SendStatus(http.StatusBadRequest)
		}
		auth, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.SendStatus(http.StatusUnauthorized)
		}
		row := db.QueryRow(`
			SELECT channel_id, user_id 
			FROM channels_users 
			WHERE channel_id = $1 
			AND user_id = $2; 
		`, channelId, auth.UserId)
		err := row.Err()
		if err != nil {
			return c.SendStatus(http.StatusForbidden)
		}
		row = db.QueryRow(`
			SELECT id, name, avatar, createdAt 
			FROM channels 
			WHERE id = $1; 
		`, channelId)
		var (
			id        string
			name      string
			avatar    *string
			createdAt string
		)
		err = row.Scan(&id, &name, &avatar, &createdAt)
		if err != nil {
			log.Println(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		if avatar == nil {
			avatar = new(string)
		}
		return c.JSON(Channel{id, name, *avatar, createdAt})
	})
}
