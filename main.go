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
// the system here use for my certification as programmer 
// the system contain http and repository as main features (domain driven i think)
// it use dependency injection with google/wire so it use object oriented programming pattern
// it use gorm as orm and postgres as database
// it use gin as http framework
// thanks to openai that help me a lot to assist me so i can finist this project faster

func main() {
	// This is only use when i need to seed the database
	seedFlag := flag.Bool("seed", false, "run database seeders and exit")
	flag.Parse()

	if *seedFlag {
		cfg := conf.NewEnvConfig()
		usersRepo := repo.NewGormUsersRepo(cfg.DB)
		profilesRepo := repo.NewGormProfilesRepo(cfg.DB)
		if err := seeder.RunAll(usersRepo, profilesRepo); err != nil {
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
