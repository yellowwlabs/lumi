package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"lumi.yellowlabs.space/internal/api/http/handlers"
	"lumi.yellowlabs.space/internal/api/http/middleware"
	"lumi.yellowlabs.space/internal/auth"
	"lumi.yellowlabs.space/internal/node"
	"lumi.yellowlabs.space/internal/organizations"
	"lumi.yellowlabs.space/internal/services"
)

func SetupRoutes(router chi.Router) {
	router.Use(middleware.CORSmiddleware)
	router.Use(middleware.SecurityHeadersMiddleware)
	router.Use(middleware.LoggingMiddleware())

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	router.Get("/health", healthHandler)

	userService := auth.NewUserService(auth.NewUserRepo(services.DB), auth.NewSessionRepo(services.DB))
	handlers.SetupAuthRoutes(router, userService)

	orgService := organizations.NewOrganizationService(organizations.NewOrganizationRepository(services.DB))
	handlers.SetupOrganizationRoutes(router, orgService)

	nodeService := node.NewNodeService(node.NewNodeRepository(services.DB), orgService)
	handlers.SetupNodeRoutes(router, nodeService)
	handlers.SetupAgentRoutes(router, nodeService)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := http.StatusOK
	checks := map[string]string{"postgres": "ok", "redis": "ok"}

	if services.DB == nil {
		checks["postgres"] = "not initialized"
		status = http.StatusServiceUnavailable
	} else if err := services.DB.PingContext(ctx); err != nil {
		checks["postgres"] = err.Error()
		status = http.StatusServiceUnavailable
	}

	if services.Redis == nil {
		checks["redis"] = "not initialized"
		status = http.StatusServiceUnavailable
	} else if err := services.Redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = err.Error()
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"status": status == http.StatusOK,
		"checks": checks,
	})
}
