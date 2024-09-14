package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

type TestPayload struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		t0 := time.Now()
		err := c.Next()
		method := c.Method()
		path := c.Path()
		status := c.Response().StatusCode()
		duration := time.Since(t0).Milliseconds()
		log.Printf("%s %s %d %dms", method, path, status, duration)
		return err
	})
	app.Use(func(c *fiber.Ctx) error {
		m := c.GetReqHeaders()
		authorizationHeaders := m["Authorization"]
		if len(authorizationHeaders) > 1 || len(authorizationHeaders) == 0 {
			return c.Next()
		}
		authorization := authorizationHeaders[0]
		token, isBearer := strings.CutPrefix(authorization, "Bearer ")
		if !isBearer || len(token) == 0 {
			return c.Next()
		}
		c.Locals("token", token)
		return c.Next()
	})
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			// c.Locals("allowed", true)
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
	app.Get("/health", func(c *fiber.Ctx) error {
		token := c.Locals("token").(string)
		log.Println("token:", token)
		return c.SendStatus(http.StatusOK)
	})
	log.Fatalln(app.Listen("127.0.0.1:8080"))
}
