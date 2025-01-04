package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

// App encapsula las dependencias de la aplicación
type App struct {
	router *gin.Engine
	logger *slog.Logger
	config Config
}

// NewApp crea una nueva instancia de la aplicación
func NewApp() *App {
	// Configuración por defecto
	config := Config{
		Port:            getEnv("PORT", "8080"),
		ShutdownTimeout: 5 * time.Second,
	}

	// Configurar logger estructurado
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Configurar router
	gin.SetMode(getEnv("GIN_MODE", "release"))
	router := gin.New()
	router.Use(
		gin.Recovery(),
		loggerMiddleware(logger),
		corsMiddleware(),
	)

	return &App{
		router: router,
		logger: logger,
		config: config,
	}
}

// setupRoutes configura las rutas de la API
func (a *App) setupRoutes() {
	a.router.GET("/health", a.healthCheck)
	a.router.GET("/uuid/:version", a.generateUUIDHandler)
}

// healthCheck maneja las verificaciones de estado
func (a *App) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

// generateUUIDHandler maneja las peticiones de generación de UUID
func (a *App) generateUUIDHandler(c *gin.Context) {
	version := c.Param("version")

	uuidString, err := a.generateUUID(version)
	if err != nil {
		var statusCode int
		if errors.Is(err, ErrUnsupportedVersion) {
			statusCode = http.StatusBadRequest
		} else {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":    uuidString,
		"version": version,
	})
}

// ErrUnsupportedVersion indica una versión de UUID no soportada
var ErrUnsupportedVersion = errors.New("versión de UUID no soportada")

// generateUUID genera un UUID del tipo especificado
func (a *App) generateUUID(version string) (string, error) {
	var uuidInstance uuid.UUID
	var err error

	switch version {
	case "0":
		uuidInstance = uuid.Nil
	case "1":
		uuidInstance, err = uuid.NewV1()
	case "4":
		uuidInstance, err = uuid.NewV4()
	case "6":
		uuidInstance, err = uuid.NewV6()
	case "7":
		uuidInstance, err = uuid.NewV7()
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedVersion, version)
	}

	if err != nil {
		return "", fmt.Errorf("error generando UUID: %w", err)
	}

	return uuidInstance.String(), nil
}

// run inicia el servidor HTTP
func (a *App) run() error {
	a.setupRoutes()

	srv := &http.Server{
		Addr:    ":" + a.config.Port,
		Handler: a.router,
	}

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)
	go func() {
		a.logger.Info("servidor iniciado", "port", a.config.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Canal para señales de sistema
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("error del servidor: %w", err)

	case <-shutdown:
		a.logger.Info("iniciando apagado graceful")

		ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("error durante el apagado: %w", err)
		}

		a.logger.Info("servidor detenido exitosamente")
		return nil
	}
}

// loggerMiddleware configura el middleware de logging
func loggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		logger.Info("request completado",
			"path", path,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
			"client_ip", c.ClientIP(),
		)
	}
}

// corsMiddleware configura el middleware CORS
func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	})
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	app := NewApp()
	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}
