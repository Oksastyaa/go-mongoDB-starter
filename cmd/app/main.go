package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"go-mongoDb-starter/database"
	"go-mongoDb-starter/pkg"
	"os"
	"os/signal"
	"time"
)

type Config struct {
}

func main() {
	// Setup logger
	pkg.SetupLogger()
	//.env
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// Config
	db, err := database.InitMongoDB()
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}

	e := echo.New()
	//	init controller
	c := NewController(db, os.Getenv("JWT_SECRET"))
	app := Config{}

	// Routes
	app.Routes(e, c, []byte(os.Getenv("JWT_SECRET")))
	// Server
	app.server(e)

}

func (app *Config) server(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		if err := e.Start(":" + port); err != nil {
			logrus.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logrus.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server shutdown")
}
