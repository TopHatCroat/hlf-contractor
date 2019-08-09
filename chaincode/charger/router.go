package main

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

func CreateRouter(router *router.Group) {
	router.
		Query("QueryById", QueryById, param.String("contractor"), param.String("charge_id")).
		Query("QueryAll", QueryAll).
		Invoke("InvokeStartChargeTransaction", InvokeStartChargeTransaction, param.Struct("startTransaction", &StartTransaction{})).
		Invoke("InvokeStopTransaction", InvokeStopChargeTransaction, param.Struct("stopTransaction", &StopTransaction{})).
		Invoke("InvokeCompleteTransaction", InvokeCompleteChargeTransaction, param.Struct("completeTransaction", &CompleteTransaction{}))
}
