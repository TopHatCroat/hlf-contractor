package fabric

import (
	"fmt"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"github.com/pkg/errors"
)

type Client struct {
	Sdk             *fabsdk.FabricSDK
	ContextProvider *context.ClientProvider
	CA              *msp.CAClientImpl
}

func New(configPath string) (*Client, error) {

	logging.SetLevel("fabsdk/fab", logging.DEBUG)
	logging.SetLevel("fabsdk/client", logging.DEBUG)

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

func (c *Client) GetChannelClient(identity *shared.Identity, channelName string) (*channel.Client, error) {
	channelCtx := c.Sdk.ChannelContext("default", fabsdk.WithUser(identity.Id), fabsdk.WithOrg("AwesomeAgency"))

	channelClient, err := channel.New(channelCtx)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error getting '%s' channel client", channelName))
	}

	return channelClient, nil
}
