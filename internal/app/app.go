package app

import (
	"context"

	"table-app/conf"
	"table-app/gui"
	"table-app/internal/log"

	"github.com/pkg/errors"
)

type Application struct {
	Gui    *gui.App
	ctx    context.Context
	logger *log.Adapter
	config conf.Remote

	cancel  context.CancelFunc
	runners []Runner
	closers []Closer
}

func New() *Application {
	logger, err := log.NewFromConfig(*logConfig())
	if err != nil {
		panic(err)
	}

	return &Application{
		ctx:    context.Background(),
		logger: logger,
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
	//select {
	//case err := <-errChan:
	//	return err
	//case <-a.ctx.Done():
	//	return nil
	//}
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
