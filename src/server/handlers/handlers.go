package handlers

import (
	"tacklo/db"
	"tacklo/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func CreateRoom(d *sqlx.DB, c *fiber.Ctx, rooms map[string]*types.Room) error {
	id := uuid.New().String()
	err := db.CreateRoom(d, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	r := types.Room{Register: make(chan *types.Player),
		Players:    make(map[int]*types.Player),
		Unregister: make(chan *types.Player),
		Broadcast:  make(chan int),
		Start:      false,
	}
	rooms[id] = &r
	go r.Run()
	return c.JSON(fiber.Map{"status": 200, "roomId": id})

}
