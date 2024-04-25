package main

import (
	"patients/internal/app"
	"patients/internal/models/patient"
	"patients/internal/models/patient/file"
	"patients/pkg/logging"
)

const configPath = "config.yml"

func main() {
	// Создаем экземпляр логгера
	logger := logging.GetLogger()

	// Получаем конфигурацию приложения из файла конфигурации
	cfg := app.GetConfig(logger, configPath)

	// Создаем экземпляр репозитория пациентов
	repo := file.NewPatientRepo(logger)

	// Создаем экземпляр хэндлера пациентов и передаем в него репозиторий и логгер
	cfg.Handler = patient.NewHandler(repo, logger)

	// Запускаем приложение с полученной конфигурацией
	err := cfg.Start(logger)
	if err != nil {
		// Если при запуске приложения произошла ошибка, записываем ее в лог и завершаем работу приложения
		logger.Fatalf("%v", err)
	}
}
