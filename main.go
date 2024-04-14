package main

import (
	"log"
	"reflect"
	"sync"

	"github.com/probablynotkai/connection"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	fc := &connection.FlatConnection{
		Directory: "C:\\Users\\kaih2\\Desktop\\Development\\AuthK",
	}

	go func() {
		defer wg.Done()
		connectToDataSource(fc)
	}()

	wg.Wait()

	// users := fc.GetUsers()
	// if len(*users) == 0 {
	// 	user, err := fc.CreateUser("CouldBeKai")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}

	// 	log.Println("No users exist, created " + user.Name)
	// }

	groups := fc.GetGroups()
	for _, v := range *groups {
		if v.Name == "admin" {
			permitted, err := v.Can("write post")
			if err != nil {
				log.Fatal(err)
				return
			}

			if permitted {
				log.Println("Post has been created by " + v.Name + " role.")
			}
		}
	}
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
