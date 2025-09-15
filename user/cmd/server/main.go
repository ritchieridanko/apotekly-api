package main

import "github.com/ritchieridanko/apotekly-api/user/internal/server"

func main() {
	app := server.New()
	app.Run()
}
