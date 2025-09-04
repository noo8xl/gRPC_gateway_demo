package utils

import "log"

// GetLocalServerAddress -> get a dev server address
func GetLocalServerAddress(serviceName string) string {
	switch serviceName {
	case "auth":
		return "0.0.0.0:10009"
	case "profile":
		return "0.0.0.0:10010"
	case "reviews":
		return "0.0.0.0:3000"
	case "orders":
		return "0.0.0.0:4000"
	case "blog":
		return "0.0.0.0:5000"
	case "notifications":
		return "0.0.0.0:6000"
	case "payments":
		return "0.0.0.0:7000"
	case "offers":
		return "0.0.0.0:8000"
	case "promotions":
		return "0.0.0.0:9000"
	case "gateway":
		return "0.0.0.0:20003"
	default:
		log.Printf("Got a wrong service name: %s", serviceName)
		return ""
	}
}

// GetProductionServerAddress -> get a prod server address
func GetProductionServerAddress(serviceName string) string {
	switch serviceName {
	case "auth":
		return ""
	case "profile":
		return ""
	case "reviews":
		return ""
	case "orders":
		return ""
	case "blog":
		return ""
	case "notifications":
		return ""
	case "payments":
		return ""
	case "offers":
		return ""
	case "promotions":
		return ""
	case "gateway":
		return ""
	default:
		log.Printf("Got a wrong service name: %s", serviceName)
		return ""
	}
}
