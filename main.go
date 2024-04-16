package main

import (
	"log"
	"reflect"
	"sync"

	"github.com/probablynotkai/connection"
)

var (
	Connection *connection.FlatConnection
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	Connection = &connection.FlatConnection{
		Directory: "",
	}

	go func() {
		defer wg.Done()
		connectToDataSource(Connection)
	}()

	wg.Wait()

	// users := Connection.GetUsers()
	// if len(*users) == 0 {
	// 	user, err := Connection.CreateUser("CouldBeKai")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}

	// 	log.Println("No users exist, created " + user.Name)
	// }

	groups := Connection.GetGroups()
	for i := range *groups {
		group := &(*groups)[i]

		if group.Name == "admin" {
			permitted, err := group.Can("write post")
			if err != nil {
				log.Fatal(err)
				return
			}

			if permitted {
				log.Println("Post has been created by " + group.Name + " role.")
			}
		}
	}

	// group := Connection.GetGroup(1)
	// if group != nil {
	// 	err := group.Grant("test2")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	err = Connection.Save()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}

func connectToDataSource(source any) {
	if source == nil {
		log.Fatal("nil data source provided")
		return
	}

	if reflect.TypeOf(source).String() == "*connection.FlatConnection" {
		dataSource := source.(*connection.FlatConnection)
		dataSource.Connect()
	} else if reflect.TypeOf(source).String() == "*connection.SQLConnection" {
		dataSource := source.(*connection.SQLConnection)
		dataSource.Connect()
	}
}
