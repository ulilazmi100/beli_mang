package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/jackc/pgx/v5/pgxpool"

	configs "beli_mang/cfg"
	local_mid "beli_mang/middleware"
)

type Server struct {
	dbPool    *pgxpool.Pool
	app       *fiber.App
	validator *validator.Validate
}

func NewServer(db *pgxpool.Pool) *Server {
	app := fiber.New(fiber.Config{
		Prefork:      false, // Disable prefork as it may increase memory usage per process (originally it's disabled by default though)
		Concurrency:  70,    // Adjust this based on testing, e.g., start with 10
		ErrorHandler: local_mid.ErrorHandler,
	})

	validate := validator.New()

	app.Use(logger.New())
	app.Use(pprof.New())

	return &Server{
		dbPool:    db,
		app:       app,
		validator: validate,
	}
}

func (s *Server) Run(config configs.Config) {
	s.app.Listen(":" + config.APPPort)
}
