package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Error: Missing arguments.")
		fmt.Println("Usage: go run add-user.go <email> <name> <password> <role>")
		fmt.Println("Example: go run add-user.go prabhatkonly@gmail.com \"Prabhat\" Hello@123 student")
		os.Exit(1)
	}

	email := os.Args[1]
	name := os.Args[2]
	password := os.Args[3]
	role := os.Args[4]

	// Use environment variable or prompt for DSN
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		fmt.Println("No POSTGRES_DSN env var set. Attempting connection...")
		fmt.Print("Enter your PostgreSQL superuser password (default user 'postgres'): ")
		var pgPass string
		fmt.Scanln(&pgPass)
		pgPass = strings.TrimSpace(pgPass)

		// Build DSN to connect to the default 'postgres' database first as superuser
		dsn = fmt.Sprintf("postgres://postgres:%s@localhost:5432/postgres?sslmode=disable", pgPass)
	}

	ctx := context.Background()

	// 1. Connect to postgres default database first
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		fmt.Printf("Error: Unable to connect to database: %v\n", err)
		fmt.Println("\nTips to solve this:")
		fmt.Println("1. Make sure your local PostgreSQL service is running on your machine.")
		fmt.Println("2. Make sure the password you entered for the 'postgres' user is correct.")
		fmt.Println("3. If you use a custom user/password/port, set the POSTGRES_DSN environment variable:")
		fmt.Println("   $env:POSTGRES_DSN=\"postgres://username:password@localhost:5432/postgres?sslmode=disable\"")
		os.Exit(1)
	}

	// 2. Automatically create the 'skillofide' database user role if it doesn't exist
	var roleExists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = 'skillofide')").Scan(&roleExists)
	if err != nil {
		fmt.Printf("Error checking for 'skillofide' role: %v\n", err)
		conn.Close(ctx)
		os.Exit(1)
	}

	if !roleExists {
		fmt.Println("Creating PostgreSQL user role 'skillofide' with password 'password'...")
		_, err = conn.Exec(ctx, "CREATE ROLE skillofide WITH LOGIN PASSWORD 'password' SUPERUSER")
		if err != nil {
			// Try without superuser if permissions are limited
			_, err = conn.Exec(ctx, "CREATE ROLE skillofide WITH LOGIN PASSWORD 'password' CREATEDB")
			if err != nil {
				fmt.Printf("Warning: Failed to create role 'skillofide': %v\n", err)
			} else {
				fmt.Println("Role 'skillofide' created successfully with CREATEDB privileges.")
			}
		} else {
			fmt.Println("Role 'skillofide' created successfully with SUPERUSER privileges.")
		}
	} else {
		fmt.Println("PostgreSQL role 'skillofide' already exists.")
	}

	// 3. Check and create 'skillofide' database if it doesn't exist
	var dbExists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'skillofide')").Scan(&dbExists)
	if err != nil {
		fmt.Printf("Error checking for 'skillofide' database: %v\n", err)
		conn.Close(ctx)
		os.Exit(1)
	}

	if !dbExists {
		fmt.Println("Database 'skillofide' does not exist. Creating it...")
		_, err = conn.Exec(ctx, "CREATE DATABASE skillofide OWNER skillofide")
		if err != nil {
			fmt.Printf("Error creating database 'skillofide': %v\n", err)
			conn.Close(ctx)
			os.Exit(1)
		}
		fmt.Println("Database 'skillofide' created successfully.")
	} else {
		fmt.Println("Database 'skillofide' already exists.")
	}
	conn.Close(ctx)

	// 4. Connect directly to the 'skillofide' database to set up table and insert user
	// Switch the target database in DSN to skillofide
	targetDsn := strings.Replace(dsn, "/postgres?", "/skillofide?", 1)
	if !strings.Contains(dsn, "/postgres?") {
		targetDsn = dsn
	}

	targetConn, err := pgx.Connect(ctx, targetDsn)
	if err != nil {
		fmt.Printf("Error: Unable to connect to 'skillofide' database: %v\n", err)
		os.Exit(1)
	}
	defer targetConn.Close(ctx)

	// 5. Ensure users table exists
	_, err = targetConn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email      TEXT NOT NULL UNIQUE,
			name       TEXT NOT NULL,
			password   TEXT NOT NULL,
			role       TEXT NOT NULL DEFAULT 'student',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		fmt.Printf("Error: Unable to create or verify users table: %v\n", err)
		os.Exit(1)
	}

	// 6. Insert or update user
	query := `
		INSERT INTO users (email, name, password, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) 
		DO UPDATE SET name = EXCLUDED.name, password = EXCLUDED.password, role = EXCLUDED.role, updated_at = now();
	`
	_, err = targetConn.Exec(ctx, query, email, name, password, role)
	if err != nil {
		fmt.Printf("Error: Failed to insert/update user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully added/updated user:\n")
	fmt.Printf("  Name:     %s\n", name)
	fmt.Printf("  Email:    %s\n", email)
	fmt.Printf("  Role:     %s\n", role)
	fmt.Println("\nYour local database is now completely configured for the SkillofIDE backend!")
}
