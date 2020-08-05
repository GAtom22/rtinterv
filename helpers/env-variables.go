package helpers

import(
	"github.com/joho/godotenv"
	"os"
	"log"
)

// GetEnvVariable environment variable from the .env file
func GetEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}