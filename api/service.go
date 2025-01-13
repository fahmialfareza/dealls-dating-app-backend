package api

import (
	"fmt"
	"os"
)

func initRunServices() (map[string]bool, string) {
	services := make(map[string]bool)

	args := os.Args[1:] // Exclude the first argument, which is the program name

	if len(args) < 1 {
		fmt.Println("Usage: go run main.go <service>")
		fmt.Println("Usage: You can run some services like go run main.go http grpc")
		panic("Please provide the service name that want you run")
	}

	for _, arg := range args {
		services[arg] = true
	}

	// get the service that run
	var kvSlice []struct {
		Key   string
		Value bool
	}
	for k, v := range services {
		kvSlice = append(kvSlice, struct {
			Key   string
			Value bool
		}{k, v})
	}

	// Access the first element of the slice
	firstKey := kvSlice[0].Key

	return services, firstKey
}
