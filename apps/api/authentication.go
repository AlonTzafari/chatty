package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthCtx struct {
	UserId string
}

type LoginReq struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

func Authentication(app *fiber.App, db *sql.DB) {
	app.Use(AuthMiddleware(db))
	app.Post("/api/login", loginHandler(db))
	app.Post("/api/register", registerHandler(db))
	app.Post("/api/logout", logoutHandler(db))
}

func AuthMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		m := c.GetReqHeaders()
		apiKey := m["X-Api-Key"]
		if len(apiKey) == 1 && apiKey[0] != "" && apiKey[0] == os.Getenv("API_KEY") {
			c.Locals("auth", AuthCtx{apiKey[0]})
			return c.Next()
		}
		session := c.Cookies("session")
		if session == "" {
			return c.Next()
		}
		rows, err := db.Query("SELECT userId FROM sessions WHERE id = $1 AND expiresAt > $2", session, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return c.Next()
		}
		isResult := rows.Next()
		if !isResult {
			c.ClearCookie("session")
			return c.Next()
		}
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			c.ClearCookie("session")
			return c.Next()
		}
		c.Locals("auth", AuthCtx{userId})
		return c.Next()
	}
}
func loginHandler(db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.ClearCookie("session")
		var loginReq LoginReq
		err := c.BodyParser(&loginReq)
		if err != nil || loginReq.Username == "" || loginReq.Password == "" {
			return c.SendStatus(http.StatusBadRequest)
		}
		log.Printf("loginReq username: %v, password: %v", loginReq.Username, loginReq.Password)
		row := db.QueryRow("SELECT id, password FROM users WHERE username = $1", loginReq.Username)
		var (
			id       string
			password string
		)
		err = row.Scan(&id, &password)
		if err != nil {
			return c.SendStatus(http.StatusUnauthorized)
		}
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginReq.Password))
		if err != nil {
			return c.SendStatus(http.StatusUnauthorized)
		}
		expiresAt := time.Now().UTC().Add(time.Millisecond * 1000 * 60 * 30)
		row = db.QueryRow(`INSERT INTO sessions (userId, expiresAt) VALUES ($1, $2) RETURNING id;`, id, expiresAt.Format(time.RFC3339))
		defer db.Exec(`DELETE FROM sessions WHERE expiresAt < $1;`, time.Now().UTC().Format(time.RFC3339))
		var sessionId string
		err = row.Scan(&sessionId)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		c.Cookie(&fiber.Cookie{Name: "session", Value: sessionId, HTTPOnly: true})
		return c.SendStatus(http.StatusOK)
	}
}
func registerHandler(db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var registerReq LoginReq
		err := c.BodyParser(&registerReq)
		if err != nil || registerReq.Username == "" || registerReq.Password == "" {
			return c.SendStatus(http.StatusBadRequest)
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), 10)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		_, err = db.Exec(`
	INSERT INTO users (username, password) 
	VALUES ($1, $2)`, registerReq.Username, string(hashed))
		if err != nil {
			if e, ok := err.(*pq.Error); ok {
				code := pq.ErrorCode("23505")
				if e.Code == code {
					return c.SendStatus(http.StatusUnauthorized)
				}
			}
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusOK)
	}
}
func logoutHandler(db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("session")
		if sessionId == "" {
			return c.SendStatus(http.StatusOK)
		}
		_, err := db.Exec(`DELETE FROM sessions WHERE id = $1`, sessionId)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		c.ClearCookie("session")
		return c.SendStatus(http.StatusOK)
	}
}
