package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/di"
	"faizalmaulana/lsp/models/repo"
	"faizalmaulana/lsp/models/seeder"
)

// Hallo This is Faizal Maulana Cashier system

func main() {
	seedFlag := flag.Bool("seed", false, "run database seeders and exit")
	flag.Parse()

	if *seedFlag {
		cfg := conf.NewEnvConfig()
		usersRepo := repo.NewGormUsersRepo(cfg.DB)
		if err := seeder.RunAll(usersRepo); err != nil {
			log.Fatalf("seeding failed: %v", err)
		}
		log.Println("seeding completed")
		return
	}

	srv := di.InitializeServer()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
