package main

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

func CreateRouter(router *router.Group) {
	router.
		Query("QueryById", QueryById, param.String("mspId"), param.String("email")).
		Query("QueryAll", QueryAll).
		Invoke("InvokeBlockUserTransaction", InvokeBlockUserTransaction, param.String("mspId"), param.String("email")).
		Invoke("InvokeUnblockUserTransaction", InvokeUnblockUserTransaction, param.String("mspId"), param.String("email"))
}
