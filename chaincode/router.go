package chaincode

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/modules/application"
	"github.com/TopHatCroat/hlf-contractor/chaincode/modules/project"
	"github.com/TopHatCroat/hlf-contractor/chaincode/schema"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	"github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = mapping.StateMappings{}.
			Add(
			&schema.Project{},
			mapping.PKeySchema(&schema.ProjectId{}),
			mapping.List(&schema.ProjectList{})).
		Add(
			&schema.Application{},
			mapping.PKeySchema(&schema.ApplicationId{}),
			mapping.List(&schema.ApplicationList{}),
		)
	EventMappings = mapping.EventMappings{}.
			Add(&schema.PublishProject{}).
			Add(&schema.PublishApplication{})
)

func CreateRouter(router *router.Group) {
	router.Use(mapping.MapStates(StateMappings))
	router.Use(mapping.MapEvents(EventMappings))

	router.Group("Project").
		Query(`List`, project.QueryList).
		Query(`Get`, project.QueryById, defparam.Proto(&schema.ProjectId{})).
		Invoke(`Publish`, project.InvokePublish, defparam.Proto(&schema.PublishProject{}))

	router.Group("Application").
		Query(`GetByContractor`, application.QueryByContractor).
		Query(`Get`, application.QueryById, defparam.Proto(&schema.ApplicationId{})).
		Invoke(`Publish`, application.InvokePublish, defparam.Proto(&schema.PublishApplication{}))
}
