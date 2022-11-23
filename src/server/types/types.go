package types

import (
	// "fmt"
	"strconv"
	"tacklo/utils"

	"github.com/gofiber/websocket/v2"
)

type Room struct {
	Players    map[int]*Player
	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan int
	State      [3][3]int8
	Start      bool
}

type Player struct {
	Conn   *websocket.Conn
	Send   chan interface{}
	Done   chan int8
	Chance bool
	Score  int8
}

func (r *Room) Run() {
	// fmt.Println("Room Started for listening")
	for {
		select {
		case player := <-r.Register:
			// fmt.Println("player registration called")
			r.Players[len(r.Players)+1] = player
			// fmt.Println("player Register")

		case player := <-r.Unregister:
			var key int
			for id, p := range r.Players {
				if p == player {
					key = id
				}
			}
			delete(r.Players, key)
			// fmt.Println("i called deleter")

		case code := <-r.Broadcast:
			// fmt.Println("inside Broadcast")
			if code == 0 {
				// fmt.Println("player Added called")
				//Here Player Added to the Room so we have to broadcast it in the room
				msg := map[string]interface{}{
					"mes":  "plyadd",
					"data": len(r.Players),
				}
				for _, p := range r.Players {
					// fmt.Println("message sending")
					p.Send <- msg
				}
			} else if code == 1 {
				// fmt.Println("player deleted called")
				msg := map[string]interface{}{
					"mes": "plydeleted",
				}
				for _, p := range r.Players {
					p.Send <- msg
				}
			} else if code == 2 {
				// fmt.Println("game started called")
				var oscore int8
				for id, p := range r.Players {
					if id == 1 {
						oscore = r.Players[2].Score
					} else {
						oscore = r.Players[1].Score
					}

					msg := map[string]interface{}{
						"mes":    "justchecking",
						"you":    p.Score,
						"other":  oscore,
						"chance": p.Chance,
						"id":     id,
					}
					p.Send <- msg
				}
			} else if code == 3 {
				// fmt.Println("game updater called")
				comp := utils.IsGameCompleted(&r.State)
				if comp == true {
					r.Start = false
					var oscore int8
					for id, p := range r.Players {
						if id == 1 {
							oscore = r.Players[2].Score
						} else {
							oscore = r.Players[1].Score
						}
						if p.Chance == true {
							p.Score += 1
						}

						msg := map[string]interface{}{
							"mes":   "gmend",
							"you":   p.Score,
							"other": oscore,
						}
						p.Send <- msg
					}
				} else {
					for _, p := range r.Players {
						p.Chance = !p.Chance
						msg := map[string]interface{}{
							"mes":     "gmud",
							"payload": r.State,
							"chance":  p.Chance,
						}
						p.Send <- msg
					}
				}
			} else if code == 4 {
				msg := map[string]interface{}{
					"mes": "badreq",
				}
				for _, p := range r.Players {
					p.Send <- msg
				}
			}

		}
	}
}

func (player *Player) Writer() {
	// fmt.Println("Player Writer Fired")
	for {
		select {
		case msg := <-player.Send:
			// fmt.Println(msg)
			player.Conn.WriteJSON(msg)
		}
	}
}

func (player *Player) Reader(rooms map[string]*Room, roomId *string) {
	var msg map[string]string
	room := rooms[*roomId]
	for {
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			delete(rooms, *roomId)
			room.Unregister <- player
			room.Broadcast <- 2
			player.Done <- 0
			// fmt.Println("read:", err)
			break
		}
		if msg["mes"] == "start" {
			if len(room.Players) == 2 {
				room.State = [3][3]int8{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
				chance, _ := strconv.Atoi(msg["payload"])
				// fmt.Println(chance)
				for id, p := range room.Players {
					if id == chance {
						p.Chance = true
					} else {
						p.Chance = false
					}
				}
				room.Start = true
				room.Broadcast <- 2
			} else {
				room.Broadcast <- 4
			}
		} else if msg["mes"] == "gmud" {
			if room.Start == true && player.Chance == true {
				row, _ := strconv.Atoi(msg["row"])
				column, _ := strconv.Atoi(msg["column"])
				if room.State[row][column] == 0 {
					room.State[row][column] = 1
					room.Broadcast <- 3
				}
			} else {
				room.Broadcast <- 4
			}
		}
	}
}
