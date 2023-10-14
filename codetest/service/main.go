package main

import (
	"service/router"
)

func main() {
	r := router.Router()

	r.Run()
}
