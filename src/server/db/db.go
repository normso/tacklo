package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDb(url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateRoom(db *sqlx.DB, rId string) error {
	_, err := db.Exec("INSERT INTO tacklo_room (uuid,players,status) VALUES($1,0,0)", rId)
	if err != nil {
		return err
	}
	return nil
}

func IsRoomExist(db *sqlx.DB, rId string) bool {
	rs := []Room{}
	db.Select(&rs, "select uuid from tacklo_room where uuid=$1", rId)
	if len(rs) == 0 {
		return false
	}
	return true
}
