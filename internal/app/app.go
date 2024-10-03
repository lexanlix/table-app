package app

import (
	"context"
	"os"
	"path"
	"strings"

	"table-app/gui"
	"table-app/internal/log"

	"github.com/pkg/errors"
)

type Application struct {
	ctx           context.Context
	logger        *log.Adapter
	Gui           *gui.App
	MigrationsDir string

	cancel  context.CancelFunc
	runners []Runner
	closers []Closer
}

func New() *Application {
	logger, err := log.NewFromConfig(*logConfig())
	if err != nil {
		panic(err)
	}

	isDev := strings.ToLower(os.Getenv("APP_MODE")) == "dev"

	migrationsDir, err := migrationsDirPath(isDev)
	if err != nil {
		logger.Fatal(context.Background(), "resolve migrations dir path")
	}

	return &Application{
		ctx:           context.Background(),
		logger:        logger,
		MigrationsDir: migrationsDir,
	}
}

func (a *Application) Context() context.Context {
	return a.ctx
}

func (a *Application) Logger() *log.Adapter {
	return a.logger
}

func (a *Application) AddRunners(runners ...Runner) {
	a.runners = append(a.runners, runners...)
}

func (a *Application) AddClosers(closers ...Closer) {
	a.closers = append(a.closers, closers...)
}

func (a *Application) Run() error {
	errChan := make(chan error)

	for i := range a.runners {
		go func(index int, runner Runner) {
			err := runner.Run(a.ctx)
			if err != nil {
				select {
				case errChan <- errors.WithMessagef(err, "start runner[%d] -> %T", index, runner):
				default:
					a.logger.Error(a.ctx, errors.WithMessagef(err, "start runner[%d] -> %T", index, runner))
				}
			}
		}(i, a.runners[i])
	}

	a.Gui.Run()
	return nil
}

func (a *Application) Shutdown() {
	a.Close()
}

func (a *Application) Close() {
	for i := 0; i < len(a.closers); i++ {
		closer := a.closers[i]
		err := closer.Close()
		if err != nil {
			a.logger.Error(a.ctx, errors.WithMessagef(err, "closers[%d] -> %T", i, closer))
		}
	}
}

func logConfig() *log.Config {
	return &log.Config{
		InitialLevel: -1,
	}
}

func migrationsDirPath(isDev bool) (string, error) {
	if isDev {
		return "./migrations", nil
	}

	return relativePathFromBin("migrations")
}

func relativePathFromBin(part string) (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", errors.WithMessage(err, "get executable path")
	}
	return path.Join(path.Dir(ex), part), nil
}
