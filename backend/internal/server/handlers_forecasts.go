package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kfc-forecast/internal/models"
)

type hourlyEntry struct {
	Hour              int     `json:"hour"`
	PredictedQuantity float64 `json:"predicted_quantity"`
}

type productForecast struct {
	Product models.Product `json:"product"`
	Hourly  []hourlyEntry  `json:"hourly"`
}

type forecastResponse struct {
	Store     models.Store      `json:"store"`
	Date      string            `json:"date"`
	Forecasts []productForecast `json:"forecasts"`
}

func (s *Server) getForecasts(c *gin.Context) {
	storeIDStr := c.Query("store_id")
	dateStr := c.Query("date")

	if storeIDStr == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "store_id and date query params are required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be in YYYY-MM-DD format"})
		return
	}

	var store models.Store
	if err := s.db.First(&store, storeIDStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch store"})
		}
		return
	}

	var forecasts []models.Forecast
	if err := s.db.
		Preload("Product").
		Where("store_id = ? AND forecast_date = ?", store.ID, date).
		Order("product_id, hour").
		Find(&forecasts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch forecasts"})
		return
	}

	c.JSON(http.StatusOK, forecastResponse{
		Store:     store,
		Date:      dateStr,
		Forecasts: groupByProduct(forecasts),
	})
}

func (s *Server) triggerGenerate(c *gin.Context) {
	if s.adminToken != "" && c.GetHeader("X-Admin-Token") != s.adminToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid X-Admin-Token header"})
		return
	}

	go func() {
		if err := s.forecastSvc.Generate(); err != nil {
			log.Printf("triggerGenerate: async generation failed: %v", err)
		}
	}()

	loc := s.forecastSvc.Loc()
	tomorrow := time.Now().In(loc).AddDate(0, 0, 1).Format("2006-01-02")
	c.JSON(http.StatusAccepted, gin.H{
		"message": "forecast generation started",
		"date":    tomorrow,
	})
}

// groupByProduct converts a flat list of forecast rows into a product-keyed structure.
func groupByProduct(forecasts []models.Forecast) []productForecast {
	seen := make(map[uint]int) // product_id → index in result
	result := []productForecast{}

	for _, f := range forecasts {
		idx, exists := seen[f.ProductID]
		if !exists {
			result = append(result, productForecast{
				Product: *f.Product,
				Hourly:  []hourlyEntry{},
			})
			idx = len(result) - 1
			seen[f.ProductID] = idx
		}
		result[idx].Hourly = append(result[idx].Hourly, hourlyEntry{
			Hour:              f.Hour,
			PredictedQuantity: f.PredictedQuantity,
		})
	}

	return result
}
