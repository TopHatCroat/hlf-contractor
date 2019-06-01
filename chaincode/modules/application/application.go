package application

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/schema"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
)

func QueryByContractor(c router.Context) (interface{}, error) {
	return c.State().List(&schema.Application{})
}

func QueryById(c router.Context) (interface{}, error) {
	id := c.Param().(*schema.ApplicationId)
	result, err := c.State().Get(id, &schema.Application{})
	return result, err
}

func InvokePublish(c router.Context) (res interface{}, err error) {
	publishData := c.Param().(*schema.PublishApplication)

	// Validate the input message defined in schema
	if err = publishData.Validate(); err != nil {
		return nil, errors.Wrap(err, "Payload invalid")
	}

	contractor, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, errors.Wrap(err, "Error obtaining identity")
	}

	projectId := schema.ProjectId{
		Issuer: publishData.ProjectIssuer,
		Name:   publishData.ProjectName,
	}

	projectResult, err := c.State().Get(projectId, &schema.Project{})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not find project with name = %s", projectId)
	}

	projectResultId := &schema.ProjectId{
		Name:   projectResult.(schema.Project).Name,
		Issuer: projectResult.(schema.Project).Issuer,
	}

	//serializedIdentity := contractor.ToSerialized().String()

	// Create state entry
	project := &schema.Application{
		Contractor:  contractor.Cert.Subject.CommonName,
		ProjectId:   projectResultId.String(),
		Price:       publishData.Price,
		State:       schema.Application_APPLIED,
		ApplyDate:   ptypes.TimestampNow(),
		Description: publishData.Description,
	}

	if err = c.Event().Set(publishData); err != nil {
		return nil, err
	}

	return project, c.State().Insert(project)
}
