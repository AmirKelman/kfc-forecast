package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kfc-forecast/internal/models"
)

func (s *Server) listStores(c *gin.Context) {
	var stores []models.Store
	if err := s.db.Order("id").Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stores"})
		return
	}
	c.JSON(http.StatusOK, stores)
}
