package main

import (
	"contree/api"
	"contree/game"
	"contree/iohandler"
	"contree/models"
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initDB := flag.Bool("d", false, "initDB to true will create schemas and add default users")
	flag.Parse()
	models.InitDB(*initDB)

	router := gin.New()

	srv := iohandler.InitSocketIo(router)
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true

	router.Use(cors.New(cfg))
	defer srv.Close()
	api.ApplyRoutes(router)
	game.InitGame()

	router.Run(":8081")
}
