package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kubernetes-sigs/kernel-module-management/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getRuntimeDir() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")

	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}

	return filepath.Join(runtimeDir, "kmm")
}

func setupRunDir(path string) error {
	mkdir := false

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not check if path %s exists: %v", path, err)
		}

		mkdir = true
	} else {
		if !fi.IsDir() {
			if err = os.Remove(path); err != nil {
				return fmt.Errorf("could not remove file %s: %v", path, err)
			}

			mkdir = true
		}
	}

	if mkdir {
		if err = os.MkdirAll(path, os.ModeDir); err != nil {
			return fmt.Errorf("could not create directory %s: %v", path, err)
		}
	}

	return nil
}

func main() {
	runtimeDir := getRuntimeDir()

	if err := setupRunDir(runtimeDir); err != nil {
		log.Fatalf("Could not create runtime directory: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewControlPlaneClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dm, err := c.GetDesiredModules(ctx, &proto.Node{Name: "toast"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Desired Modules: %s", dm)
}
