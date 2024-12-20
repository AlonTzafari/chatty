package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Id        string
	ChannelId string
	UserId    string
	Content   string
	CreatedAt string
}

type MessageDisplay struct {
	Message
	Username string
}

type CreateMessageReq struct {
	Content   string `json:"content" xml:"content" form:"content"`
	ChannelId string `json:"channelId" xml:"channelId" form:"channelId"`
}

func Messages(app *fiber.App, db *sql.DB) {
	app.Post("/api/messages", func(c *fiber.Ctx) error {
		authCtx, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.SendStatus(http.StatusUnauthorized)
		}
		var createMessageReq CreateMessageReq
		err := c.BodyParser(&createMessageReq)
		if err != nil || createMessageReq.ChannelId == "" || createMessageReq.Content == "" {
			log.Print(err)
			return c.SendStatus(http.StatusBadRequest)
		}
		isInChannel, err := isUserInChannel(authCtx.UserId, createMessageReq.ChannelId, db)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		if !isInChannel {
			return c.SendStatus(http.StatusForbidden)
		}

		now := time.Now().UTC().Format(time.RFC3339)
		row := db.QueryRow(
			`INSERT INTO messages (channel_id, user_id, content, createdAt) VALUES ($1, $2, $3, $4) RETURNING id;`,
			createMessageReq.ChannelId,
			authCtx.UserId,
			createMessageReq.Content,
			now,
		)
		var messageId string
		err = row.Scan(&messageId)
		if err != nil {
			log.Print("MESSAGE INSERT", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		msg, err := json.Marshal(Message{
			messageId,
			createMessageReq.ChannelId,
			authCtx.UserId,
			createMessageReq.Content,
			now,
		})
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		defer PubsubClient.Send("message-updates", msg)
		return c.Status(http.StatusCreated).JSON(struct {
			Id string `json:"id"`
		}{messageId})

	})
	app.Get("/api/messages", func(c *fiber.Ctx) error {
		channelId := c.Query("channel_id")
		if channelId == "" {
			return c.SendStatus(http.StatusBadRequest)
		}
		auth, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.SendStatus(http.StatusUnauthorized)
		}
		isInChannel, err := isUserInChannel(auth.UserId, channelId, db)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		if !isInChannel {
			return c.SendStatus(http.StatusForbidden)
		}

		rows, err := db.Query(`
			SELECT messages.id, messages.user_id, messages.content, messages.createdAt, users.username
			FROM messages
			LEFT JOIN users
			ON messages.user_id = users.id
			WHERE channel_id = $1;
		`, channelId)

		if err != nil {
			log.Println(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		var messages []MessageDisplay = []MessageDisplay{}
		for rows.Next() {
			var (
				id        string
				userId    string
				content   string
				createdAt string
				username  string
			)
			err := rows.Scan(&id, &userId, &content, &createdAt, &username)
			if err != nil {
				log.Println(err)
				return c.SendStatus(http.StatusInternalServerError)
			}
			messages = append(messages, MessageDisplay{Message{id, channelId, userId, content, createdAt}, username})
		}
		return c.JSON(messages)
	})
}

func isUserInChannel(userId string, channelId string, db *sql.DB) (bool, error) {
	row := db.QueryRow(`
			SELECT 1
			FROM channels_users
			WHERE channel_id = $1
			AND user_id = $2;`,
		channelId,
		userId,
	)
	err := row.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
