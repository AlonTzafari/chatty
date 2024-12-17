package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

type User struct {
	Id       string
	Username string
}

func main() {
	dbUri := os.Getenv("DB_URI")
	port := os.Getenv("PORT")
	address := os.Getenv("ADDRESS")
	if dbUri == "" {
		dbUri = "postgresql://root:password@localhost:5432/app?sslmode=disable"
	}
	if port == "" {
		port = "8080"
	}
	if address == "" {
		address = "127.0.0.1"
	}

	db, err := sql.Open("postgres", dbUri)
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
	log.Fatalln(app.Listen(fmt.Sprintf("%v:%v", address, port)))
}
