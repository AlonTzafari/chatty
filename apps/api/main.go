package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

type WSMessage struct {
	Channel string                 `json:"channel"`
	Payload map[string]interface{} `json:"payload"`
}

func decodePayload[T any](m map[string]interface{}, out *T) error {
	str, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(str, out)
}

type User struct {
	Id       string
	Username string
}
type TestPayload struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func main() {
	connStr := "postgresql://root:password@localhost:5432/app?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(cors.New())
	app.Use(LoggerMiddleware)
	Authentication(app, db)
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			}
			log.Printf("mt: %d, recv: %s", mt, msg)
			var m WSMessage
			if err = json.Unmarshal(msg, &m); err != nil {
				break
			}
			var testPayload TestPayload
			if err = decodePayload(m.Payload, &testPayload); err != nil {
				break
			}
			log.Println("testPayload:", testPayload)
			if err = c.WriteMessage(mt, msg); err != nil {
				break
			}
		}
	}))

	app.Get("/users", func(c *fiber.Ctx) error {
		page := c.QueryInt("page")
		rows, err := db.Query("SELECT id, username FROM users LIMIT 5 OFFSET $1", 5*page)
		if err != nil {
			log.Fatalln(err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		defer rows.Close()
		var users []User
		for rows.Next() {
			var (
				id       string
				username string
			)
			err := rows.Scan(&id, &username)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			user := User{id, username}
			users = append(users, user)
		}
		if users == nil {
			return c.SendStatus(404)
		}
		return c.JSON(users)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		token := c.Locals("token").(string)
		log.Println("token:", token)
		return c.SendStatus(http.StatusOK)
	})
	log.Fatalln(app.Listen("127.0.0.1:8080"))
}
