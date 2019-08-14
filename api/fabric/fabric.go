package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
)

type Client struct {
	Sdk             *fabsdk.FabricSDK
	ContextProvider *context.ClientProvider
	CA              *msp.CAClientImpl
}

func New(configPath string) (*Client, error) {
	sdkConfig := config.FromFile(configPath)
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		fmt.Printf("failed to create new SDK: %s\n", err)
	}

	sdkCtxProvider := sdk.Context(fabsdk.WithOrg("AwesomeAgency"))
	sdkCtx, err := sdkCtxProvider()
	if err != nil {
		fmt.Printf("failed to create new SDK: %s\n", err)
	}

	caClient, err := msp.NewCAClient("AwesomeAgency", sdkCtx)

	return &Client{
		Sdk:             sdk,
		ContextProvider: &sdkCtxProvider,
		CA:              caClient,
	}, nil
}

func (c *Client) Register(username, password string) error {
	req := &api.RegistrationRequest{
		Name:   username,
		Secret: password,
	}

	_, err := c.CA.Register(req)
	return err
}

func (c *Client) Login(username, password string) error {
	req := &api.EnrollmentRequest{
		Name:   username,
		Secret: password,
	}

	err := c.CA.Enroll(req)
	return err
}
