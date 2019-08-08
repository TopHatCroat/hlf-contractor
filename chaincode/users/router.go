package main

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/modules/user"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/router/param/defparam"
)

func CreateRouter(router *router.Group) {
	router.
		Group("User").
		Query(`QueryById`, user.QueryById, param.String("email")).
		Query(`QueryAll`, user.QueryAll)
}
