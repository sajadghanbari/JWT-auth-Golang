package main

import (
	"JWT-Authentication-go/config"
	db "JWT-Authentication-go/data/database"
)

func main() {
	var cfg = config.GetConfig()
	db.InitDb(cfg)
	defer
	db.CloseDb()
}