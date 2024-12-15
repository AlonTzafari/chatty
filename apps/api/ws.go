package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WSMessage struct {
	Channel string                 `json:"channel"`
	Payload map[string]interface{} `json:"payload"`
}

type SubscribePayload struct {
	Topic string `json:"topic"`
}

func WS(app *fiber.App) {
	app.Use("/api/ws", func(c *fiber.Ctx) error {
		if _, ok := c.Locals("auth").(AuthCtx); !ok {
			return fiber.ErrUnauthorized
		}
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	},
		websocket.New(func(wsConn *websocket.Conn) {
			logger := log.Default()
			defer wsConn.Close()
			userId := "anon"
			if user, ok := wsConn.Locals("auth").(AuthCtx); ok {
				userId = user.UserId
			}
			logger.Println("user", userId, "connected")
			var (
				msg []byte
				err error
			)
			subscriber := PubsubClient.AddSubscriber()
			logger.Println("user", userId, "PubsubClient.AddSubscriber()")
			defer subscriber.Close()
			channel := subscriber.GetChannel()
			logger.Println("user", userId, "subscriber.GetChannel()")
			go WSPublish(wsConn, channel, func(m PubsubMsg) bool {
				return true
			}, logger, userId)
			logger.Println("user", userId, "go WSPublish()")
			for {
				if _, msg, err = wsConn.ReadMessage(); err != nil {
					logger.Println("wsConn.ReadMessage()", err)
					break
				}
				logger.Println("user", userId, string(msg))
				var m WSMessage
				if err = json.Unmarshal(msg, &m); err != nil {
					logger.Println("json.Unmarshal(msg, &m)", err)
					break
				}
				if m.Channel == "subscribe" {
					var subPayload SubscribePayload
					if err = decodePayload(m.Payload, &subPayload); err != nil {
						logger.Println(err)
						break
					}
					subscriber.Subscribe(subPayload.Topic)
					logger.Printf("user %v subscriber.Subscribe(%v)\n", userId, subPayload.Topic)

				} else if m.Channel == "quit" {
					var subPayload SubscribePayload
					if err = decodePayload(m.Payload, &subPayload); err != nil {
						logger.Println(err)
						break
					}
					subscriber.Unsubscribe(subPayload.Topic)
					logger.Printf("user %v subscriber.Unsubscribe(%v)\n", userId, subPayload.Topic)
				}
			}
			logger.Printf("user %v close connection\n", userId)
		}, websocket.Config{Subprotocols: []string{"chat-ws"}}))
}

func WSPublish(conn *websocket.Conn, channel chan PubsubMsg, filter func(PubsubMsg) bool, logger *log.Logger, userId string) {
	for m := range channel {
		logger.Printf("received message for user %v on topic %v, payload %v", userId, m.Topic, string(m.Message))
		msg, err := SerializePubsubMsg(m)
		if err != nil {
			conn.CloseHandler()(websocket.CloseInternalServerErr, "")
			return
		}
		if err = conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func SerializePubsubMsg(m PubsubMsg) ([]byte, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(m.Message, &payload); err != nil {
		return []byte{}, err
	}
	wsMsg := WSMessage{m.Topic, payload}
	return json.Marshal(wsMsg)
}

func decodePayload[T any](m map[string]interface{}, out *T) error {
	str, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(str, out)
}
