package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func generateUUID() {
	u2 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u2)
}
