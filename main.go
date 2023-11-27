package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	controllersystem "github.com/MrDweller/controller-system/controller-system"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	address := os.Getenv("ADDRESS")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Panic(err)
	}
	systemName := os.Getenv("SYSTEM_NAME")

	serviceRegistryAddress := os.Getenv("SERVICE_REGISTRY_ADDRESS")
	serviceRegistryPort, err := strconv.Atoi(os.Getenv("SERVICE_REGISTRY_PORT"))
	if err != nil {
		log.Panic(err)
	}

	controllerSystem, err := controllersystem.NewControllerSystem(address, port, systemName, serviceRegistryAddress, serviceRegistryPort)
	if err != nil {
		log.Panic(err)
	}
	controllerSystem.StartControllerSystem()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	log.Printf("Stopping the controller system!")

	err = controllerSystem.StopControllerSystem()
	if err != nil {
		log.Panic(err)
	}
}
