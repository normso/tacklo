package routers

import (
	"fmt"
	"tacklo/handlers"
	"tacklo/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/jmoiron/sqlx"
)

func CreateRouters(app *fiber.App, db *sqlx.DB, rooms map[string]*types.Room) {

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"help": "new"})
	})

	//Router For Creating Room --> no payload required
	app.Get("/createRoom", func(c *fiber.Ctx) error {
		return handlers.CreateRoom(db, c, rooms)

	})

	//Router for connecting to per room specific websocket
	app.Get("/ws/:rid", websocket.New(func(c *websocket.Conn) {
		fmt.Println(rooms)
		rid := c.Params("rid")
		room := rooms[rid]
		player := types.Player{Conn: c,
			Send:   make(chan interface{}),
			Done:   make(chan int8),
			Chance: false, Score: 0,
		}
		room.Register <- &player
		room.Broadcast <- 0
		go player.Reader(rooms, &rid)
		go player.Writer()
		<-player.Done
		return

	}))
}
