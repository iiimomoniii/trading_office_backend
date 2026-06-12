package route

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"trading-office/trading_office_backend/config"
	_ "trading-office/trading_office_backend/docs"
	"trading-office/trading_office_backend/handler"
	"trading-office/trading_office_backend/middleware"
	"trading-office/trading_office_backend/utils"
)

type APIServer struct {
	server *fiber.App
	cfg    *config.Config
}

func NewAPIServer(cfg *config.Config, sqlDB *sql.DB) App {
	app := fiber.New(fiber.Config{
		AppName:        cfg.App.Name,
		ReadTimeout:    time.Millisecond * time.Duration(cfg.Fiber.ReadTimeout),
		WriteTimeout:   time.Millisecond * time.Duration(cfg.Fiber.WriteTimeout),
		IdleTimeout:    time.Millisecond * time.Duration(cfg.Fiber.IdleTimeout),
		ReadBufferSize: cfg.Fiber.ReadBufferSize,
		BodyLimit:      cfg.Fiber.BodyLimitSize,
	})

	app.Use(middleware.Cors())
	app.Use(middleware.Logger())

	// ─── Handlers ──────────────────────────────────────────
	dashboardHandler := handler.NewDashboardHandler()

	registerRoutes(app, dashboardHandler)

	return &APIServer{server: app, cfg: cfg}
}

func registerRoutes(
	app *fiber.App,
	dashboardHandler *handler.DashboardHandler,
) {
	// ─── Public ────────────────────────────────────────────
	app.Get("/live", dashboardHandler.Live)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Post("/auth/token",
		limiter.New(limiter.Config{
			Max:        5,
			Expiration: 1 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return utils.TooManyRequests(c)
			},
		}),
	)

	// ─── Protected Routes ──────────────────────────────────
	v1 := app.Group("/api/v1", middleware.Auth())
	_ = v1 // will be populated in upcoming phases
}

func (s *APIServer) Start() {
	addr := fmt.Sprintf(":%d", s.cfg.Fiber.Port)

	if err := s.server.Listen(addr); err != nil {
		panic(err)
	}
}

func (s *APIServer) Stop() {
	if err := s.server.Shutdown(); err != nil {
		panic(err)
	}
}
