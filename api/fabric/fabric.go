package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
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

func (c *Client) GetChannelClient(identity *api.IdentityResponse, channelName string) (*channel.Client, error) {
	channelCtx := c.Sdk.ChannelContext("default", fabsdk.WithUser(identity.ID), fabsdk.WithOrg("AwesomeAgency"))

	channelClient, err := channel.New(channelCtx)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error getting '%s' channel client", channelName))
	}

	return channelClient, nil
}

func (c *Client) AllUsers(identity *api.IdentityResponse) ([]User, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	res, err := channelClient.Query(channel.Request{
		ChaincodeID: "users",
		Fcn:         "QueryAll",
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from users::QueryAll function")
	}

	user := make([]User, 20)
	err = json.Unmarshal(res.Payload, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from users::QueryById function")
	}

	return user, nil
}

func (c *Client) FindUserById(identity *api.IdentityResponse, userName string) (*User, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	args := [][]byte{[]byte(identity.ID), []byte(identity.ID)}
	res, err := channelClient.Query(channel.Request{
		ChaincodeID: "users",
		Fcn:         "QueryById",
		Args:        args,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from users::QueryById function")
	}

	user := &User{}
	err = json.Unmarshal(res.Payload, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from users::QueryById function")
	}

	return user, nil
}
