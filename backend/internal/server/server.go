package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kfc-forecast/internal/forecast"
)

type Server struct {
	router      *gin.Engine
	db          *gorm.DB
	forecastSvc *forecast.Service
	adminToken  string
}

func New(db *gorm.DB, forecastSvc *forecast.Service, adminToken string) *Server {
	s := &Server{
		router:      gin.Default(),
		db:          db,
		forecastSvc: forecastSvc,
		adminToken:  adminToken,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.router.GET("/health", s.health)

	api := s.router.Group("/api")
	api.GET("/stores", s.listStores)
	api.GET("/forecasts", s.getForecasts)
	api.POST("/forecasts/generate", s.triggerGenerate)
}

func (s *Server) health(c *gin.Context) {
	sqlDB, err := s.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HTTPServer returns a configured *http.Server for graceful shutdown support.
func (s *Server) HTTPServer(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}
}
