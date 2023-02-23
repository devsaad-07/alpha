package main

import (
	"alpha/communication"
	"alpha/db"
)

func init() {
	//feediscount.Init_db()
	server := communication.GetServer()
	server.Run()
}

func main() {
	db.RunMigration()
}
