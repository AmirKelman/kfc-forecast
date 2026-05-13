package forecast

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"kfc-forecast/internal/models"
)

// Service generates daily sales forecasts using an average of historical sales.
type Service struct {
	db          *gorm.DB
	historyDays int
	daysAhead   int
	loc         *time.Location
	mu          sync.Mutex // prevents concurrent generation runs
}

func NewService(db *gorm.DB, historyDays, daysAhead int, loc *time.Location) *Service {
	if loc == nil {
		loc = time.UTC
	}
	return &Service{
		db:          db,
		historyDays: historyDays,
		daysAhead:   daysAhead,
		loc:         loc,
	}
}

// avgRow holds one row returned by the AVG aggregation query.
type avgRow struct {
	StoreID           uint    `gorm:"column:store_id"`
	ProductID         uint    `gorm:"column:product_id"`
	Hour              int     `gorm:"column:hour"`
	PredictedQuantity float64 `gorm:"column:predicted_quantity"`
}

// Generate computes average-based forecasts for the target date and persists them.
// It is idempotent and concurrency-safe: at most one run executes at a time.
func (s *Service) Generate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	targetDate := s.today().AddDate(0, 0, s.daysAhead)
	cutoff := s.today().AddDate(0, 0, -s.historyDays)

	log.Printf("forecast: generating for %s using %d days of history (since %s)",
		targetDate.Format("2006-01-02"), s.historyDays, cutoff.Format("2006-01-02"))

	var rows []avgRow
	err := s.db.Raw(`
		SELECT
			store_id,
			product_id,
			hour,
			ROUND(AVG(quantity)::numeric, 2) AS predicted_quantity
		FROM historical_sales
		WHERE sale_date >= ?
		  AND sale_date <  ?
		GROUP BY store_id, product_id, hour
		ORDER BY store_id, product_id, hour
	`, cutoff, s.today()).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("forecast avg query: %w", err)
	}

	if len(rows) == 0 {
		log.Println("forecast: no historical data found — skipping")
		return nil
	}

	generatedAt := time.Now()
	records := make([]models.Forecast, len(rows))
	for i, r := range rows {
		records[i] = models.Forecast{
			StoreID:           r.StoreID,
			ProductID:         r.ProductID,
			ForecastDate:      targetDate,
			Hour:              r.Hour,
			PredictedQuantity: r.PredictedQuantity,
			GeneratedAt:       generatedAt,
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete any previously generated forecasts for this date so the run is idempotent.
		if err := tx.Where("forecast_date = ?", targetDate).Delete(&models.Forecast{}).Error; err != nil {
			return fmt.Errorf("forecast clear: %w", err)
		}
		if err := tx.CreateInBatches(records, 500).Error; err != nil {
			return fmt.Errorf("forecast insert: %w", err)
		}
		log.Printf("forecast: inserted %d records for %s", len(records), targetDate.Format("2006-01-02"))
		return nil
	})
}

// Loc returns the timezone the service uses for day boundaries.
func (s *Service) Loc() *time.Location { return s.loc }

// today returns midnight of the current day in the configured timezone.
func (s *Service) today() time.Time {
	t := time.Now().In(s.loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, s.loc)
}
