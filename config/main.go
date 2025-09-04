package config

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/noo8xl/anvil-common/exceptions"
)

func init() {

	// load .env from multiple possible locations
	envFiles := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
		"../../../../.env",
	}

	loaded := false
	for _, file := range envFiles {
		if err := loadEnvFile(file); err != nil {
			if err.Error() == "no such file or directory" {
				continue
			}
		}
		loaded = true
		break
	}

	if !loaded {
		log.Printf("Warning: .env file not found in any of the search paths")
	}

}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return errors.New("no such file or directory")
		}
		return exceptions.HandleAnException(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}
