package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"log"

	"github.com/coddemn/base-structure-go/internal/config"
	"github.com/coddemn/base-structure-go/internal/database"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	ctx := context.Background()

	var (
		dir   = flag.String("dir", "migrations", "директория миграций")
		cmd   = flag.String("cmd", "up", "команда: up, down, status, redo, reset")
		steps = flag.Int("steps", 1, "сколько миграций откатить (для down)")
	)
	flag.Parse()

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set dialect: %v", err)
	}

	switch *cmd {
	case "up":
		err = goose.Up(db, *dir)
	case "down":
		if *steps <= 1 {
			err = goose.Down(db, *dir)
		} else {
			for i := 0; i < *steps; i++ {
				if err = goose.Down(db, *dir); err != nil {
					if errors.Is(err, goose.ErrNoCurrentVersion) {
						err = nil // достигли нуля
					}
					break
				}
			}
		}
	case "status":
		err = goose.Status(db, *dir)
	case "redo":
		err = goose.Redo(db, *dir)
	case "reset":
		err = goose.Reset(db, *dir)
	default:
		log.Fatalf("unknown command: %s", *cmd)
	}

	if err != nil {
		log.Fatalf("migrate %s: %v", *cmd, err)
	}

	log.Printf("✔ migrate %s — done\n", *cmd)
}
