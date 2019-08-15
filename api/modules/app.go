package modules

import (
	"github.com/TopHatCroat/hlf-contractor/api/fabric"
)

type App struct {
	Client   *fabric.Client
	sessions map[string]string
}

func NewApp(fabricConfig string) (*App, error) {
	fabClient, err := fabric.New(fabricConfig)
	if err != nil {
		return nil, err
	}

	app := &App{
		Client:   fabClient,
		sessions: make(map[string]string),
	}

	return app, nil
}

func (app *App) GetSession(token string) string {
	return "username1@mail.com"
	return app.sessions[token]
}
