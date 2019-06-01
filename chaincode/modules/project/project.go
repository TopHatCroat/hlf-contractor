package project

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/schema"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
)

func QueryList(c router.Context) (interface{}, error) {
	return c.State().List(&schema.Project{})
}

func QueryById(c router.Context) (interface{}, error) {
	id := c.Param().(*schema.ProjectId)
	result, err := c.State().Get(id, &schema.Project{})
	return result, err
}

func InvokePublish(c router.Context) (res interface{}, err error) {
	publishData := c.Param().(*schema.PublishProject)

	// Validate the input message defined in schema
	if err = publishData.Validate(); err != nil {
		return nil, errors.Wrap(err, "Payload invalid")
	}

	issuer, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, errors.Wrap(err, "Error obtaining identity")
	}

	// Create state entry
	project := &schema.Project{
		Issuer:         issuer.Cert.Subject.CommonName,
		Assessor:       publishData.Assessor,
		Name:           publishData.Name,
		OpenDate:       ptypes.TimestampNow(),
		StartDate:      publishData.StartDate,
		EndDate:        publishData.EndDate,
		EstimatedValue: publishData.EstimatedValue,
		State:          schema.Project_OPEN, // Initial state
		Description:    publishData.Description,
	}

	if err = c.Event().Set(publishData); err != nil {
		return nil, err
	}

	return project, c.State().Insert(project)
}
