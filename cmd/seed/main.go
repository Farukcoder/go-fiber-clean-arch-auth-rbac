package main

import (
	"fmt"
	"os"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/database/seed/seeders"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/config"
)

func main() {
	config.InitLogger()
	fmt.Println("INFO  Seeding database.")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("ERROR Failed to load config:", err.Error())
		os.Exit(1)
	}

	fmt.Print("Database\\Seeders\\DatabaseSeeder ................................. RUNNING")
	if err := seeders.SeedRBAC(cfg); err != nil {
		fmt.Println()
		fmt.Println("Database\\Seeders\\DatabaseSeeder ................................... 100ms FAIL")
		fmt.Println("ERROR", err.Error())
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println("Database\\Seeders\\DatabaseSeeder ................................... 15ms DONE")
	fmt.Println()

	fmt.Print("Database\\Seeders\\UserSeeder ................................. RUNNING")
	if err := seeders.SeedUsers(cfg); err != nil {
		fmt.Println()
		fmt.Println("Database\\Seeders\\UserSeeder ................................... 100ms FAIL")
		fmt.Println("ERROR", err.Error())
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println("Database\\Seeders\\UserSeeder ................................... 10ms DONE")
}
