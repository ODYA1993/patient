package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"net/http"
	"os"
	"os/signal"
	"patients/internal/models/patient"
	"patients/pkg/logging"
	"sync"
	"syscall"
	"time"
)

type App struct {
	Port    string `yaml:"port" env:"PORT" env-default:"8082"`
	IsDebug *bool  `yaml:"is_debug" env:"IS_DEBUG" env-default:"false"`
	Handler *patient.Handler
}

var cfg *App
var once sync.Once

func GetConfig(logger *logging.Logger, configPath string) *App {
	once.Do(func() {
		logger.Info("read application configuration")
		cfg = &App{}
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			logger.Fatal("configuration file not found")
		}
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Info(help)
			logger.Fatal(err)
		}

	})
	return cfg
}

func (c *App) Start(logger *logging.Logger) error {
	logger.Infof("starting server on port: %s", c.Port)

	logger.Info("configure router")
	router := c.ConfigureRouter()

	server := &http.Server{
		Addr:         ":" + c.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Создаем канал для получения сигналов завершения работы приложения
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("server stopped with error: %v", err)
		}
	}()

	// Ждем сигнала завершения работы приложения
	<-shutdownChan

	// Начинаем завершение работы сервера
	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown with error: %v", err)
	}

	logger.Info("application stopped")
	return nil
}

func (c *App) ConfigureRouter() *gin.Engine {
	router := gin.New()
	router.POST("/patient", c.Handler.NewPatient)
	router.GET("/patients", c.Handler.GetListPatients)
	router.POST("/patient/edit", c.Handler.EditPatient)
	router.POST("/patient/delete", c.Handler.DelPatient)
	return router
}
