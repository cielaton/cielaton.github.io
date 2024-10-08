package main

import (
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"server/database"
	"server/handlers"
)

func databaseMiddleware(database *pgx.Conn) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("database", database)
			return next(c)
		}
	}
}

func main() {
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Create the echo server
	echoServer := echo.New()

	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	echoServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, loggerValue middleware.RequestLoggerValues) error {
			log.Info().Str("URI", loggerValue.URI).Str("Method", loggerValue.Method).Dur("Latency", loggerValue.Latency).Str("Remote IP", loggerValue.RemoteIP).Int("Status", loggerValue.Status).Msg("request")
			return nil
		}}))

	// Connect to the database
	postgreSQL := database.ConnectDB()
	echoServer.Use(databaseMiddleware(postgreSQL))

	echoServer.GET("/development-tools", handlers.GetDevelopmentTools)
	echoServer.GET("/personal-projects", handlers.GetPersonalProjects)
	echoServer.GET("/team-projects", handlers.GetTeamProjects)

	echoServer.POST("/personal-projects", handlers.InsertPersonalProjects)
	echoServer.POST("/development-tools", handlers.InsertDevelopmentTools)

	echoServer.DELETE("/personal-projects/:id", handlers.DeletePersonalProjects)
	echoServer.DELETE("/development-tools/:id", handlers.DeleteDevelopmentTools)

	echoServer.Logger.Fatal(echoServer.Start(":1323"))
	defer func() {
		database.DisconnectDB(postgreSQL)
	}()
}
