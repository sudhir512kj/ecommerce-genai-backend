package repository

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/sudhir512kj/ecommerce_backend/config"
	"github.com/sudhir512kj/ecommerce_backend/database"
)

var dbTest database.Database
var conf *config.Config

func generateRandomEmail() string {
	rand.Seed(time.Now().UnixNano())

	// Define the valid characters for the local part and domain part
	localChars := "abcdefghijklmnopqrstuvwxyz0123456789"
	domainChars := "abcdefghijklmnopqrstuvwxyz"

	// Generate the local part
	localPartLen := rand.Intn(10) + 5 // Random length between 5 and 14 characters
	localPart := make([]byte, localPartLen)
	for i := range localPart {
		localPart[i] = localChars[rand.Intn(len(localChars))]
	}

	// Generate the domain part
	domainPartLen := rand.Intn(3) + 2 // Random length between 2 and 4 characters
	domainPart := make([]byte, domainPartLen)
	for i := range domainPart {
		domainPart[i] = domainChars[rand.Intn(len(domainChars))]
	}

	// Combine the local part and domain part to form the email address
	emailAddress := fmt.Sprintf("%s@%s.com", localPart, domainPart)
	return emailAddress
}

func TestMain(m *testing.M) {
	// Set up the test environment
	// Set up the test environment
	conf = config.GetConfig()
	dbTest = database.NewPostgresDatabase(conf)

	// Run the tests
	os.Exit(m.Run())
}
