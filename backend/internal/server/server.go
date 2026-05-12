package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
}

func New(db *gorm.DB) *Server {
	s := &Server{
		router: gin.Default(),
		db:     db,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.router.GET("/health", s.health)

	api := s.router.Group("/api")
	_ = api // handlers added in later stages
}

func (s *Server) health(c *gin.Context) {
	sqlDB, err := s.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}
