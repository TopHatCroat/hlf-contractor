package modules

import (
	"github.com/TopHatCroat/hlf-contractor/api/fabric"
)

type App struct {
	client   *fabric.Client
	sessions map[string]string
}

func NewApp(fabricConfig string) (*App, error) {
	fabClient, err := fabric.New(fabricConfig)
	if err != nil {
		return nil, err
	}

	return &App{
		client:   fabClient,
		sessions: make(map[string]string),
	}, nil
}
