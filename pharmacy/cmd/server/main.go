package main

import "github.com/ritchieridanko/apotekly-api/pharmacy/internal/server"

func main() {
	app := server.New()
	app.Run()
}
