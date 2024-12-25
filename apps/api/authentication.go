package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Id        string
	UserId    string
	Username  string
	ExpiresAt time.Time
}
type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]Session
}

var sm = &SessionManager{sessions: make(map[string]Session)}

func (sm *SessionManager) GetSession(sessionId string) (Session, bool) {
	sm.mu.Lock()
	defer func() {
		now := time.Now()
		for k, v := range sm.sessions {
			if now.After(v.ExpiresAt) {
				delete(sm.sessions, k)
			}
		}
		sm.mu.Unlock()
	}()
	session, ok := sm.sessions[sessionId]
	if !ok {
		return session, ok
	}
	now := time.Now()
	if now.After(session.ExpiresAt) {
		return Session{}, false
	}
	return session, true
}
func (sm *SessionManager) AddSession(userId string, username string) (Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sessionId := ""
	for i := 0; i < 3; i++ {
		uid, err := uuid.NewRandom()
		if err != nil {
			return Session{}, err
		}
		sessionId = uid.String()
		_, ok := sm.sessions[sessionId]
		if ok {
			continue
		}
		session := Session{Id: sessionId, UserId: userId, Username: username, ExpiresAt: time.Now().UTC().Add(time.Hour * 24)}
		sm.sessions[sessionId] = session
		return session, nil
	}
	return Session{}, fmt.Errorf("failed to add session")
}
func (sm *SessionManager) RemoveSession(sessionId string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionId)
}

type AuthCtx struct {
	UserId   string
	Username string
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
	app.Get("/api/.me", meHandler(db))
}

func AuthMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		m := c.GetReqHeaders()
		apiKey := m["X-Api-Key"]
		if len(apiKey) == 1 && apiKey[0] != "" && apiKey[0] == os.Getenv("API_KEY") {
			c.Locals("auth", AuthCtx{apiKey[0], "API"})
			return c.Next()
		}
		sessionId := c.Cookies("session")
		if sessionId == "" {
			return c.Next()
		}
		session, ok := sm.GetSession(sessionId)
		if !ok {
			c.ClearCookie("session")
			return c.Next()
		}
		c.Locals("auth", AuthCtx{session.UserId, session.Username})
		return c.Next()
	}
}
func loginHandler(db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.ClearCookie("session")
		var loginReq LoginReq
		err := c.BodyParser(&loginReq)
		if err != nil || loginReq.Username == "" || loginReq.Password == "" {
			log.Printf("Error parsing body: %v, body: %v", err, string(c.BodyRaw()))
			return c.SendStatus(http.StatusBadRequest)
		}
		log.Printf("loginReq username: %v, password: %v", loginReq.Username, loginReq.Password)
		row := db.QueryRowContext(c.Context(), "SELECT id, password FROM users WHERE username = $1", loginReq.Username)
		var (
			id       string
			password string
		)
		err = row.Scan(&id, &password)
		if err == sql.ErrNoRows {
			log.Printf("Error sql.ErrNoRows")
			return c.SendStatus(http.StatusUnauthorized)
		}
		if err != nil {
			log.Printf("Error QueryRow select user: %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginReq.Password))
		if err != nil {
			log.Printf("Error bcrypt.CompareHashAndPassword: %v", err)
			return c.SendStatus(http.StatusUnauthorized)
		}
		session, err := sm.AddSession(id, loginReq.Username)
		if err != nil {
			log.Printf("Error sm.AddSession(%v): %v", id, err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		c.Cookie(&fiber.Cookie{Name: "session", Value: session.Id, HTTPOnly: true, Expires: session.ExpiresAt})
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
		ctx, cancel := context.WithTimeout(c.Context(), 5000*time.Millisecond)
		defer cancel()
		_, err = db.ExecContext(ctx, `
	INSERT INTO users (username, password) 
	VALUES ($1, $2)`, registerReq.Username, string(hashed))
		if err != nil {
			if e, ok := err.(*pq.Error); ok {
				code := pq.ErrorCode("23505")
				if e.Code == code {
					return c.Status(http.StatusBadRequest).JSON(struct {
						Username string
						Password string
					}{"username taken", ""})
				}
			}
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusOK)
	}
}
func logoutHandler(_ *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("session")
		if sessionId == "" {
			return c.SendStatus(http.StatusOK)
		}
		sm.RemoveSession(sessionId)
		c.ClearCookie("session")
		return c.SendStatus(http.StatusOK)
	}
}
func meHandler(db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		auth, ok := c.Locals("auth").(AuthCtx)
		if !ok {
			return c.JSON(nil)
		}
		row := db.QueryRowContext(c.Context(), `SELECT username FROM users WHERE id = $1`, auth.UserId)
		var username string
		err := row.Scan(&username)
		if err != nil {
			return c.JSON(nil)
		}
		return c.JSON(User{Id: auth.UserId, Username: username})
	}
}
