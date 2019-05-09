package chaincode

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/schema"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	"github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = mapping.StateMappings{}.Add(
		&schema.Project{},
		mapping.PKeySchema(&schema.ProjectId{}),
		mapping.List(&schema.ProjectList{}),
	)
	EventMappings = mapping.EventMappings{}.
			Add(&schema.PublishProject{})
)

func NewCC() *router.Chaincode {
	r := router.New(`project`)

	// Mappings for chaincode state
	r.Use(mapping.MapStates(StateMappings))
	// Mappings for chaincode events
	r.Use(mapping.MapEvents(EventMappings))

	// Store info about the chaincode initiator
	r.Init(owner.InvokeSetFromCreator)

	// method for debug chaincode state
	debug.AddHandlers(r, `debug`, owner.Only)

	// Read all method
	r.Query(`list`, queryProjectList).
		// Get Project method, takes 2 params: Issuer Id and Project Name
		Query(`get`, queryProject, defparam.Proto(&schema.ProjectId{})).
		// txn methods
		Invoke(`publish`, invokeProjectPublish, defparam.Proto(&schema.PublishProject{}))

	return router.NewChaincode(r)
}

func queryProjectList(c router.Context) (interface{}, error) {
	return c.State().List(&schema.Project{})
}

func queryProject(c router.Context) (interface{}, error) {
	id := c.Param().(*schema.ProjectId)
	return c.State().Get(id)
}

func invokeProjectPublish(c router.Context) (res interface{}, err error) {
	publishData := c.Param().(*schema.PublishProject)

	// Validate the input message defined in schema
	if err = publishData.Validate(); err != nil {
		return nil, errors.Wrap(err, "Payload invalid")
	}

	issuer, err := identity.FromStub(c.Stub())

	// Create state entry
	issuerName := issuer.Cert.Subject.CommonName
	project := &schema.Project{
		Issuer:         issuerName,
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
