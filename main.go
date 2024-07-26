package main

import (
	"aplikasieoq/database"
	"aplikasieoq/router"
)

func main() {
	database.StartDB()

	r := router.StartApp()

	r.Run(":8081")
}
