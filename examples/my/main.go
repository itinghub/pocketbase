package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

type User struct {
	Id     string `db:"id" json:"id"`
	Status bool   `db:"status" json:"status"`
	Age    int    `db:"age" json:"age"`
	// Roles  types.JsonArray `db:"roles" json:"roles"`
}

func main() {
	app := pocketbase.New()
	app.Bootstrap()
	users := []User{}

	_ = app.Dao().DB().
		NewQuery("SELECT id, status, age, roles FROM users LIMIT 100").
		All(&users)
	log.Println(users)
}
