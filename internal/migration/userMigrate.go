package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sudhir512kj/ecommerce_backend/config"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type postgresDatabase struct {
	Db *sql.DB
}

var once sync.Once
var conf *config.Config

var db *postgresDatabase

func main() {

	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.SSLMode,
			conf.Db.TimeZone,
		)

		dbI, err := sql.Open("postgres", dsn)
		if err != nil {
			panic("failed to connect database")
		}

		db = &postgresDatabase{Db: dbI}
	})
	defer db.Db.Close()

	// Create the users table
	createUsersTable(db.Db)

	// Insert 10 demo users
	insertDemoUsers(db.Db)

	// Insert demo addresses for some users
	insertDemoAddresses(db.Db)

	fmt.Println("User migration completed successfully.")
}

func createUsersTable(db *sql.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            permissions TEXT[] DEFAULT '{buyer}',
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
}

func insertDemoUsers(db *sql.DB) {
	for i := 1; i <= 10; i++ {
		firstName := fmt.Sprintf("User%d", i)
		lastName := "Demo"
		email := fmt.Sprintf("user%d@example.com", i)
		password := "password123"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}

		permissions := []models.Permission{models.PermissionBuyer}
		if i%2 == 0 {
			permissions = append(permissions, models.PermissionSeller)
		}
		if i%3 == 0 {
			permissions = append(permissions, models.PermissionAdmin)
		}

		_, err = db.Exec(
			"INSERT INTO users (id, email, password, first_name, last_name, permissions, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			uuid.New(), email, string(hashedPassword), firstName, lastName, fmt.Sprintf("{%s}", formatPermissions(permissions)), time.Now(), time.Now(),
		)
		if err != nil {
			log.Fatalf("Error inserting demo user: %v", err)
		}
	}
}

func insertDemoAddresses(db *sql.DB) {
	userIDs := []int{1, 3, 5, 7, 9}
	for _, userID := range userIDs {
		address := &models.Address{
			UserID:  userID,
			Street:  fmt.Sprintf("Street %d", userID),
			City:    fmt.Sprintf("City %d", userID),
			State:   fmt.Sprintf("State %d", userID),
			Country: "USA",
			Zipcode: fmt.Sprintf("%05d", userID*1000),
		}

		_, err := db.Exec(
			"INSERT INTO addresses (user_id, street, city, state, country, zipcode, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			address.UserID, address.Street, address.City, address.State, address.Country, address.Zipcode, time.Now(), time.Now(),
		)
		if err != nil {
			log.Fatalf("Error inserting demo address: %v", err)
		}
	}
}

func formatPermissions(permissions []models.Permission) string {
	var formattedPermissions []string
	for _, p := range permissions {
		formattedPermissions = append(formattedPermissions, fmt.Sprintf(`"%s"`, p))
	}
	return strings.Join(formattedPermissions, ",")
}
