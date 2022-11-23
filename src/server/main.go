package main

import (
	"fmt"

	"tacklo/routers"

	"tacklo/db"
	"tacklo/types"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func main() {
	//Loading the Environment files
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error in LOading Environment file")
		return
	}

	//Connecting to Database(here it is postgres)
	d, err := db.ConnectDb(os.Getenv("db_url"))
	if err != nil {
		fmt.Println(err)
	}

	//Room initialised
	rooms := map[string]*types.Room{}

	//App initialing
	app := fiber.New()

	//middleware are here

	//First middleware CORS setup
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	//WebSocket Upgradation middleware
	app.Use("/ws/:rid", func(c *fiber.Ctx) error {
		// fmt.Println("i Called from ")
		e := db.IsRoomExist(d, c.Params("rid"))
		// fmt.Println(e)
		if e == true {
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("allowed", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		}
		return c.SendStatus(400)

	})

	// fmt.Println(db.CheckForRoom(d, "f78bd07e-0bc6-4bc2-883d-7b41f2cc6626"))

	//Creating Routers
	routers.CreateRouters(app, d, rooms)

	app.Listen(":3000")
}
