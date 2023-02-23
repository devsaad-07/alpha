package main

import (
	"alpha/communication"
	feediscount "alpha/feeDiscount"
)

func init() {
	feediscount.Init_db()
	server := communication.GetServer()
	server.Run()
}

func main() {

}
