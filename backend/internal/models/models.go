package models

import "time"

type Store struct {
	ID        uint      `gorm:"primaryKey"            json:"id"`
	Name      string    `gorm:"size:100;not null"     json:"name"`
	Location  string    `gorm:"size:200;not null"     json:"location"`
	CreatedAt time.Time `                             json:"created_at"`
}

type Product struct {
	ID        uint      `gorm:"primaryKey"                      json:"id"`
	Name      string    `gorm:"size:100;not null"               json:"name"`
	Category  string    `gorm:"size:50;not null"                json:"category"`
	Price     float64   `gorm:"type:numeric(10,2);not null"     json:"price"`
	CreatedAt time.Time `                                        json:"created_at"`
}

type HistoricalSale struct {
	ID        uint      `gorm:"primaryKey"          json:"id"`
	StoreID   uint      `gorm:"not null"            json:"store_id"`
	ProductID uint      `gorm:"not null"            json:"product_id"`
	SaleDate  time.Time `gorm:"type:date;not null"  json:"sale_date"`
	Hour      int       `gorm:"not null"            json:"hour"`
	Quantity  int       `gorm:"not null"            json:"quantity"`
	CreatedAt time.Time `                           json:"created_at"`
}

func (HistoricalSale) TableName() string { return "historical_sales" }

type Forecast struct {
	ID                uint      `gorm:"primaryKey"                  json:"id"`
	StoreID           uint      `gorm:"not null"                    json:"store_id"`
	ProductID         uint      `gorm:"not null"                    json:"product_id"`
	ForecastDate      time.Time `gorm:"type:date;not null"          json:"forecast_date"`
	Hour              int       `gorm:"not null"                    json:"hour"`
	PredictedQuantity float64   `gorm:"type:numeric(10,2);not null" json:"predicted_quantity"`
	GeneratedAt       time.Time `                                   json:"generated_at"`

	Store   *Store   `gorm:"foreignKey:StoreID"   json:"store,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
