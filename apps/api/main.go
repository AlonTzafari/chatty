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
