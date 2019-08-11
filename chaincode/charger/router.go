package main

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/charger/charge"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

func CreateRouter(router *router.Group) {
	router.Init(Init, param.String("initialPrice")).
		Query("QueryById", QueryById, param.String("contractor"), param.String("chargeId")).
		Query("QueryAll", QueryAll).
		Invoke("InvokeStartChargeTransaction", InvokeStartChargeTransaction, param.Struct("startTransaction", &charge.StartTransaction{})).
		Invoke("InvokeStopTransaction", InvokeStopChargeTransaction, param.Struct("stopTransaction", &charge.StopTransaction{})).
		Invoke("InvokeCompleteTransaction", InvokeCompleteChargeTransaction, param.Struct("completeTransaction", &charge.CompleteTransaction{}))
}
