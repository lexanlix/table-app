package assembly

import (
	"context"
	"encoding/json"

	"table-app/conf"
	"table-app/gui"
	"table-app/internal/app"
	db "table-app/internal/db/client"
	"table-app/internal/log"

	"github.com/pkg/errors"
)

type Assembly struct {
	logger       *log.Adapter
	db           *db.Client
	shutdownFunc func()
}

func New(app *app.Application) *Assembly {
	logger := app.Logger()
	dbCli := db.NewClient(logger)

	return &Assembly{
		logger:       logger,
		db:           dbCli,
		shutdownFunc: app.Shutdown,
	}
}

func (a *Assembly) ReceiveConfig(ctx context.Context, remoteConfig []byte) (*gui.App, error) {
	var newCfg conf.Remote

	err := a.UpgradeConfig(remoteConfig, &newCfg)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "upgrade remote config"))
	}

	err = a.db.Upgrade(ctx, newCfg.Database)
	if err != nil {
		return nil, errors.WithMessage(err, "upgrade db client")
	}

	locator := NewLocator(a.db, a.logger)

	// создание данных для gui с последующим занесением куда-то в ран или еще куда
	guiApp, err := locator.Config(ctx, newCfg, a.shutdownFunc)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "get locator config"))
	}

	return guiApp, nil
}

func (a *Assembly) Runners() []app.Runner {
	return []app.Runner{
		app.RunnerFunc(func(ctx context.Context) error {
			return nil
		}),
	}
}

func (a *Assembly) Closers() []app.Closer {
	return []app.Closer{
		app.CloserFunc(func() error {
			return nil
		}),
		a.db,
	}
}

func (a *Assembly) UpgradeConfig(newCfg []byte, config *conf.Remote) error {
	if len(newCfg) == 0 {
		return errors.New("new config is empty")
	}

	err := json.Unmarshal(newCfg, &config)
	if err != nil {
		return errors.WithMessage(err, "unmarshal new config")
	}

	return nil
}
