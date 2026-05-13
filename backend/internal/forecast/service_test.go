package forecast_test

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kfc-forecast/internal/forecast"
	"kfc-forecast/internal/models"
)

// setupDB spins up a throwaway Postgres container and runs AutoMigrate.
// The container is terminated automatically when the test finishes.
func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()

	pg, err := tcpg.Run(ctx, "postgres:16-alpine",
		tcpg.WithDatabase("testdb"),
		tcpg.WithUsername("test"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pg.Terminate(ctx) })

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Store{},
		&models.Product{},
		&models.HistoricalSale{},
		&models.Forecast{},
	); err != nil {
		t.Fatalf("auto-migrate: %v", err)
	}

	return db
}

func seedSale(db *gorm.DB, storeID, productID uint, hour, qty, daysAgo int) {
	db.Create(&models.HistoricalSale{
		StoreID:   storeID,
		ProductID: productID,
		SaleDate:  time.Now().UTC().AddDate(0, 0, -daysAgo),
		Hour:      hour,
		Quantity:  qty,
	})
}

// TestGenerate_AverageIsCorrect verifies that the AVG calculation matches
// the expected value for a set of known historical quantities.
func TestGenerate_AverageIsCorrect(t *testing.T) {
	db := setupDB(t)

	cases := []struct {
		name       string
		quantities []int
		wantAvg    float64
	}{
		{"equal_quantities", []int{10, 10, 10}, 10.0},
		{"varying_quantities", []int{10, 12, 14}, 12.0},
		{"single_day", []int{7}, 7.0},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			store := models.Store{Name: tc.name, Location: "Test City"}
			db.Create(&store)
			product := models.Product{Name: tc.name, Category: "Burgers", Price: 9.99}
			db.Create(&product)

			for i, qty := range tc.quantities {
				seedSale(db, store.ID, product.ID, 12, qty, i+1)
			}

			svc := forecast.NewService(db, 30, 1, time.UTC)
			if err := svc.Generate(); err != nil {
				t.Fatalf("Generate() error: %v", err)
			}

			var result models.Forecast
			if err := db.Where("store_id = ? AND product_id = ? AND hour = 12",
				store.ID, product.ID).First(&result).Error; err != nil {
				t.Fatalf("no forecast found: %v", err)
			}
			if result.PredictedQuantity != tc.wantAvg {
				t.Errorf("predicted_quantity = %.2f, want %.2f",
					result.PredictedQuantity, tc.wantAvg)
			}
		})
	}
}

// TestGenerate_Idempotent verifies that calling Generate() twice does not
// produce duplicate forecast rows.
func TestGenerate_Idempotent(t *testing.T) {
	db := setupDB(t)

	store := models.Store{Name: "Idempotent Store", Location: "Test City"}
	db.Create(&store)
	product := models.Product{Name: "Test Burger", Category: "Burgers", Price: 9.99}
	db.Create(&product)
	seedSale(db, store.ID, product.ID, 10, 8, 1)

	svc := forecast.NewService(db, 30, 1, time.UTC)
	_ = svc.Generate()
	_ = svc.Generate() // second run must replace, not append

	var count int64
	db.Model(&models.Forecast{}).Count(&count)
	if count != 1 {
		t.Errorf("forecast row count = %d after 2 Generate() calls, want 1", count)
	}
}
