package main

import (
	"alpha/communication"
)

func init() {
	//feediscount.Init_db()
	// db.RunMigration()
	server := communication.GetServer()
	server.Run()
}

func main() {
}
