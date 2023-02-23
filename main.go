package main

import (
	"fiber_websocket/db"
	"fiber_websocket/routers"
	"fiber_websocket/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	utils.LoadEnv()
	db.InitDB()
	routers.InitRouter(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", utils.ApiPort)))
}
