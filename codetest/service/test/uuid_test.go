package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestUuid(t *testing.T) {
	UUid := uuid.NewV4().String()
	fmt.Println(UUid)
}
