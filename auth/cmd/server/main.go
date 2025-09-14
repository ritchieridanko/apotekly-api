package main

import "github.com/ritchieridanko/apotekly-api/auth/internal/server"

func main() {
	app := server.New()
	app.Run()
}
