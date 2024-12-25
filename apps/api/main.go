package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(5)
	go printDBStats(db, context.Background())
	app := fiber.New(fiber.Config{WriteTimeout: 5 * time.Second})
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

func printDBStats(db *sql.DB, ctx context.Context) {
	for {
		stats := db.Stats()
		if stats.WaitCount > 0 {
			log.Printf("DB STATS %+v\n", stats)
		}
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(5 * time.Second)
		}
	}
}
