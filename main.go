package main

import (
	"alpha/communication"
)

func init() {
	server := communication.GetServer()
	server.Run()
}

func main() {

}
