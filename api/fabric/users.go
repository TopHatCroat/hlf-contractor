package fabric

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
	"github.com/pkg/errors"
)

const (
	defaultAffiliation = "awesome"
)

type User struct {
	MSP     string `json:"msp,omitempty"`
	Email   string `json:"email,omitempty"`
	Balance int    `json:"balance,omitempty"`
	State   string `json:"state,omitempty"`
}

func (c *Client) Register(username, password string, role shared.Role) error {
	req := &api.RegistrationRequest{
		Name:   username,
		Secret: password,
		//Affiliation: defaultAffiliation,
		Attributes: []api.Attribute{{Name: "role", Value: string(role), ECert: true}},
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
func (c *Client) AllUsers(identity *shared.Identity) ([]User, error) {
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

func (c *Client) FindUserById(identity *shared.Identity, userName string) (*User, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	args := [][]byte{[]byte(identity.Id), []byte(identity.Id)}
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
