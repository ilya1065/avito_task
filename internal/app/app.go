package app

import (
	"avito_task/config"
	v1 "avito_task/internal/controller/http/v1"
	"avito_task/internal/controller/http/v1/hendler"
	repo2 "avito_task/internal/repo"
	"avito_task/internal/service"
	"avito_task/pkg/postgres"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	log.Println("[APP] Starting service…")

	cfg := config.MustLoad()
	log.Printf("[CONFIG] Loaded config: HTTP_PORT=%s PGDSN=%s\n", cfg.HTTPPort, cfg.PGDSN)

	pg, err := postgres.NewPostgres(cfg.PGDSN)
	if err != nil {
		log.Fatalf("[POSTGRES] Failed to connect: %v", err)
	}
	log.Println("[POSTGRES] Connected successfully")
	defer func() {
		_ = pg.Close()
		log.Println("[POSTGRES] Connection closed")
	}()

	repo := repo2.NewRepositories(pg)
	log.Println("[REPO] Repositories initialized")

	newService := service.NewService(repo)
	log.Println("[SERVICE] Services initialized")

	h := hendler.NewHendler(newService)
	router := v1.NewRouter(h)
	log.Println("[ROUTER] HTTP router initialized")

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("[SERVER] Starting HTTP server on %s…\n", addr)

	if err = http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("[SERVER] Failed to start: %v", err)
	}
}
