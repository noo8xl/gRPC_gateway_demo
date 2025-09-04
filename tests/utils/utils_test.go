package utils_test

import (
	"log"
	"strings"
	"testing"

	utils "github.com/noo8xl/anvil-gateway/utils/server"
)

var (
	localAddresses = map[string]string{
		"auth":          "0.0.0.0:10009",
		"profile":       "0.0.0.0:10010",
		"reviews":       "0.0.0.0:3000",
		"orders":        "0.0.0.0:4000",
		"blog":          "0.0.0.0:5000",
		"notifications": "0.0.0.0:6000",
		"payments":      "0.0.0.0:7000",
		"offers":        "0.0.0.0:8000",
		"promotions":    "0.0.0.0:9000",
		"gateway":       "0.0.0.0:20003",
		// "wrong":         "Got a wrong service name: wrong",
	}
)

func TestGetLocalServerAddress(t *testing.T) {
	for serviceName, address := range localAddresses {
		log.Printf("serviceName: %s, address: %s", serviceName, address)
		result := utils.GetLocalServerAddress(serviceName)
		if strings.Contains(result, "Got a wrong service name:") {
			continue
		}
		if result != address {
			t.Errorf("TestGetLocalServerAddress error: address is %s, expected `%s`", result, address)
		}
	}
}
