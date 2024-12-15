package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

type User struct {
	Id       string
	Username string
}

func main() {
	connStr := "postgresql://root:password@localhost:5432/app?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := fiber.New()
	app.Use(cors.New())
	app.Use(LoggerMiddleware)
	Authentication(app, db)
	Channels(app, db)
	Users(app, db)
	Messages(app, db)
	WS(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	log.Fatalln(app.Listen("127.0.0.1:8080"))
}
