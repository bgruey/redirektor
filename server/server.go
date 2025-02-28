package main

import (
	"redirektor/server/api"
	"redirektor/server/model"
	"redirektor/server/repo"
)

func main() {

	migrate()

	handler := api.NewAPIHandler(8080)

	handler.Run()
}

func migrate() {
	db := repo.NewPostgresClient()

	db.DB.AutoMigrate(&model.Redirect{})
	db.DB.AutoMigrate(&model.ApiKey{})
}
