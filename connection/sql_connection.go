package connection

import (
	"log"
)

type SQLConnection struct {
	ConnectionString string
}

func (s SQLConnection) Connect() {
	log.Println("coming soon")
}
