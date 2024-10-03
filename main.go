package main

import (
	"os"

	"table-app/assembly"
	"table-app/internal/app"
	"table-app/internal/shutdown"

	"github.com/pkg/errors"
)

func main() {
	app := app.New()
	logger := app.Logger()

	remoteCfg, err := readConfigFromFile()
	if err != nil {
		logger.Fatal(app.Context(), err.Error())
	}

	// запуск БД, создание кэша и подключение к gui хендлеров
	assembly := assembly.New(app)

	app.Gui, err = assembly.ApplyConfig(app.Context(), remoteCfg)
	if err != nil {
		logger.Fatal(app.Context(), errors.WithMessage(err, "failed to apply config"))
	}

	// добавление функций старта и завершения
	app.AddRunners(assembly.Runners()...)
	app.AddClosers(assembly.Closers()...)

	// подключение точки выхода из приложения
	shutdown.On(func() {
		logger.Info(app.Context(), "starting shutdown")
		app.Gui.Shutdown()
		logger.Info(app.Context(), "shutdown completed")
	})

	// старт приложения
	err = app.Run()
	if err != nil {
		app.Shutdown()
		logger.Fatal(app.Context(), err)
	}
}

func readConfigFromFile() ([]byte, error) {
	bytes, err := os.ReadFile("conf/app_config.json")
	if err != nil {
		return nil, errors.WithMessage(err, "read app config")
	}

	return bytes, nil
}
