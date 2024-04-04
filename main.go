package main

import (
	"log"

	"github.com/probablynotkai/connection"
)

var (
	ConnectionChannel = make(chan bool)
)

func main() {
	fc := &connection.FlatConnection{
		Directory: "C:\\Users\\kharrison\\Documents\\AuthK",
	}

	go fc.Connect(ConnectionChannel)

	_ = <-ConnectionChannel

	users := fc.GetUsers()
	if len(*users) == 0 {
		user, err := fc.CreateUser("CouldBeKai")
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("No users exist, created " + user.Name)
	}
}

// func connectToDataSource(source any) {
// 	if source == nil {
// 		log.Fatal("nil data source provided")
// 		return
// 	}

// 	if reflect.TypeOf(source).Name() == "FlatConnection" {
// 		dataSource := source.(connection.FlatConnection)
// 		dataSource.Connect()
// 	} else if reflect.TypeOf(source).Name() == "SQLConnection" {
// 		dataSource := source.(connection.SQLConnection)
// 		dataSource.Connect()
// 	}
// }
