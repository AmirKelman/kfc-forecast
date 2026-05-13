package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"kfc-forecast/internal/forecast"
	"kfc-forecast/internal/models"
)

type Scheduler struct {
	cron        *cron.Cron
	forecastSvc *forecast.Service
	db          *gorm.DB
	loc         *time.Location
}

func New(forecastSvc *forecast.Service, db *gorm.DB, loc *time.Location) *Scheduler {
	if loc == nil {
		loc = time.UTC
	}
	return &Scheduler{
		cron:        cron.New(cron.WithLocation(loc)),
		forecastSvc: forecastSvc,
		db:          db,
		loc:         loc,
	}
}

// Start registers the cron job and — if no forecast exists for tomorrow — runs immediately.
func (s *Scheduler) Start(cronExpr string) error {
	_, err := s.cron.AddFunc(cronExpr, s.run)
	if err != nil {
		return err
	}
	s.cron.Start()
	log.Printf("scheduler: started with cron %q in timezone %s", cronExpr, s.loc)

	s.runIfNeeded()
	return nil
}

// Stop gracefully shuts down the cron runner.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("scheduler: stopped")
}

// run executes one forecast generation cycle and logs the outcome.
func (s *Scheduler) run() {
	start := time.Now()
	log.Printf("scheduler: forecast generation started at %s", start.Format(time.RFC3339))

	if err := s.forecastSvc.Generate(); err != nil {
		log.Printf("scheduler: forecast generation FAILED (%v)", err)
		return
	}

	log.Printf("scheduler: forecast generation completed in %s", time.Since(start).Round(time.Millisecond))
}

// runIfNeeded generates forecasts immediately on startup if tomorrow has none yet.
func (s *Scheduler) runIfNeeded() {
	tomorrow := s.tomorrowDate()

	var count int64
	s.db.Model(&models.Forecast{}).
		Where("forecast_date = ?", tomorrow).
		Count(&count)

	if count == 0 {
		log.Println("scheduler: no forecasts found for tomorrow — running startup generation")
		s.run()
	} else {
		log.Printf("scheduler: %d forecasts already exist for %s — skipping startup run",
			count, tomorrow.Format("2006-01-02"))
	}
}

func (s *Scheduler) tomorrowDate() time.Time {
	t := time.Now().In(s.loc)
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, s.loc)
}
