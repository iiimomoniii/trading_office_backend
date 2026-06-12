package route

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"trading-office/trading_office_backend/config"
	"trading-office/trading_office_backend/db"
)

type App interface {
	Start()
	Stop()
}

func Bootstrap() (App, *sync.WaitGroup, error) {

	// 1. โหลด config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("[bootstrap] config error: %w", err)
	}

	// 2. รัน migration ด้วย connection แยก
	if err := db.RunMigrations(cfg.Database); err != nil {
		fmt.Printf("[bootstrap] Migration warning: %v\n", err)
	}

	// 3. Connect database สำหรับ app
	sqlDB, err := db.NewDatabase(cfg.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("[bootstrap] ❌ DB connect failed: %w", err)
	}
	fmt.Println("[bootstrap] ✅ DB ready")

	// 4. สร้าง server
	apiServer := NewAPIServer(cfg, sqlDB)

	// 5. graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	addShutdownHook(&wg, sqlDB, func() {
		apiServer.Stop()
	})

	return apiServer, &wg, nil
}

func addShutdownHook(wg *sync.WaitGroup, sqlDB *sql.DB, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer wg.Done()
		<-c
		fmt.Println("\nShutting down...")
		f()
		if sqlDB != nil {
			sqlDB.Close()
			fmt.Println("[bootstrap] DB closed")
		}
	}()
}
