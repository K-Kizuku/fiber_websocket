package routers

import (
	"fiber_websocket/cruds"
	"fiber_websocket/db"
	"fiber_websocket/utils"
	"fiber_websocket/ws"
	"fmt"

	"strings"

	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitRouter(app *fiber.App) {
	app.Use(cors.New())
	api := app.Group("/api/v1")
	rooms := make(map[*ws.Room]bool)
	// room := ws.NewRoom()

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!!!!!!!!!!")
	})

	api.Get("/test", func(c *fiber.Ctx) error {
		temp := db.UserInfo{Id: "1"}
		result2 := db.Db.First(&temp)
		fmt.Println(result2)
		if errors.Is(result2.Error, gorm.ErrRecordNotFound) {
			log.Fatal(result2.Error)
		}
		return c.JSON(&temp)
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	api.Post(("/gamestart"), func(c *fiber.Ctx) error {
		payload := c.Body()
		slices := strings.Split(string(payload), " ")
		u, _ := uuid.NewRandom()
		store := db.Room{Room: u.String(), Id1: slices[1], Id2: slices[3], IsStreaming: true}
		db.Db.Select("room", "id1", "id2", "isStreaming").Create(&store)
		// db.Db.Select("id","name").Create()
		// rooms[c.Body()]=true
		rooms[ws.NewRoom()] = true
		return c.SendString(u.String())
	})

	app.Post("/temp", func(c *fiber.Ctx) error {

		// if err := c.BodyParser(p); err != nil {
		//     return err
		// }

		// log.Println(p.Player1) // john
		// log.Println(p.Player2) // doe

		return c.Send(c.Body())
		// ...
	})

	// app.Get("/ws/:id", room.ServeWs())
	// app.Get("/ws/:id", room.ServeWs())

	api.Post("/signUp", func(c *fiber.Ctx) error {
		var payload db.SignUpUser
		err := c.BodyParser(&payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		u, err := cruds.CreateUser(payload.Name, payload.Email, payload.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		c.JSON(&u)
		return nil
	})
}
func middleware(c *fiber.Ctx) {
	authorizationHeader := c.GetRespHeader("Authorization")
	if authorizationHeader != "" {
		ary := strings.Split(authorizationHeader, " ")
		if len(ary) == 2 {
			if ary[0] == "Bearer" {
				t, err := jwt.ParseWithClaims(ary[1], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
					return utils.SigningKey, nil
				})

				if claims, ok := t.Claims.(*jwt.MapClaims); ok && t.Valid {
					userId := (*claims)["sub"].(string)
					c.Set("user_id", userId)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
	c.Next()
}
