package main

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/router/param/defparam"
)

func CreateRouter(router *router.Group) {
	router.Group("Charge").
		Query("QueryById", QueryById, param.String("contractor"), param.String("charge_id")).
		Query("QueryAll", QueryAll).
		Invoke("InvokeStartTransaction", InvokeStartTransaction, defparam.Proto(&StartTransaction{})).
		Invoke("InvokeStopTransaction", InvokeStopChargeTransaction, defparam.Proto(&StopTransaction{})).
		Invoke("InvokeCompleteTransaction", InvokeCompleteChargeTransaction, defparam.Proto(&CompleteTransaction{}))
}
