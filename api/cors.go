package api

import (
	"github.com/go-chi/cors"
)

func corsOption(corsAllowLocalhost bool) cors.Options {
	allowedOrigins := []string{"https://*"}
	if corsAllowLocalhost {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	}
	return cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           7200, // Maximum value not ignored by any of major browsers
	}
}
