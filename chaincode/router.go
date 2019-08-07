package chaincode

import (
	"github.com/TopHatCroat/hlf-contractor/chaincode/modules/charge"
	"github.com/TopHatCroat/hlf-contractor/chaincode/modules/user"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/router/param/defparam"
)

func CreateRouter(router *router.Group) {
	router.
		Group("Charge").
		Query(`QueryById`, charge.QueryById, param.String("contractor"), param.String("charge_id")).
		Query(`QueryAll`, charge.QueryAll).
		Invoke(`InvokeStartTransaction`, charge.InvokeStartTransaction, defparam.Proto(&charge.StartTransaction{})).
		Invoke(`InvokeStopTransaction`, charge.InvokeStopChargeTransaction, defparam.Proto(&charge.StopTransaction{})).
		Invoke(`InvokeCompleteTransaction`, charge.InvokeCompleteChargeTransaction, defparam.Proto(&charge.CompleteTransaction{})).
		Group("User").
		Query(`QueryById`, user.QueryById, param.String("email")).
		Query(`QueryAll`, user.QueryAll)
}
