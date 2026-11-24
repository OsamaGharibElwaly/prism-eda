package main

import (
	"log"

	"prismeda/internal/api"
	"prismeda/internal/config"
	"prismeda/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	svc := service.New(cfg.MaxUploadBytes, cfg.SamplePath)

	gin.SetMode(gin.ReleaseMode)

	router := api.NewRouter(svc, cfg.WebDir)
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	log.Printf("Prism EDA %s — http://localhost:%s", api.Version, cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
