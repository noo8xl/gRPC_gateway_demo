package middlewares

import (
	"net/http"
	"strings"

	utils "github.com/noo8xl/anvil-gateway/utils/loggers"
)

// Logger is a middleware that logs the request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		loggerDto := utils.LoggerDto{
			Ip:        r.RemoteAddr,
			Method:    r.Method,
			Url:       r.URL.Path,
			UserAgent: r.UserAgent(),
			Referer:   r.Referer(),
			Host:      r.Host,
		}

		utils.RequestLoggerUtil(loggerDto)
		next.ServeHTTP(w, r)
	})
}
