package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	t0 := time.Now()
	err := c.Next()
	method := c.Method()
	path := c.Path()
	status := c.Response().StatusCode()
	duration := time.Since(t0).Milliseconds()
	log.Printf("%s %s %d %dms", method, path, status, duration)
	return err
}
