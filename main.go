package main

import (
	"log"
	"reflect"

	"github.com/probablynotkai/connection"
)

func main() {
	myConnection := connection.FlatConnection{
		FileLocation: "data.json",
	}

	connectToDataSource(myConnection)

	user := &myConnection.GetUsers()[0]
	myConnection.Grant(user, "test_permission")
}

func connectToDataSource(source any) {
	if source == nil {
		log.Fatal("nil data source provided")
		return
	}

	if reflect.TypeOf(source).Name() == "FlatConnection" {
		dataSource := source.(connection.FlatConnection)
		dataSource.Connect()
	} else if reflect.TypeOf(source).Name() == "SQLConnection" {
		dataSource := source.(connection.SQLConnection)
		dataSource.Connect()
	}
}
