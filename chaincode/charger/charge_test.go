package main

import (
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"io/ioutil"
	"path"
	"testing"

	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	contractorName = "contractor"

	globalMSP = "globalMSP"
	localMSP  = "localMSP"

	cryptoPath      = path.Join("testdata")
	globalAdminCert = path.Join(cryptoPath, "global", "admin", "cert.pem")
	globalUserCert  = path.Join(cryptoPath, "global", "user", "cert.pem")
	localAdminCert  = path.Join(cryptoPath, "local", "admin", "cert.pem")
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Charge Suite")
}

var _ = Describe("Charge", func() {
	r := router.New("root")
	debug.AddHandlers(r, "debug", owner.Only)
	CreateRouter(r)
	r.Init(owner.InvokeSetFromCreator)
	cc := testcc.NewMockStub("charge", router.NewChaincode(r))

	// Load actor certificates
	globalActors, err := testcc.IdentitiesFromFiles(
		globalMSP,
		map[string]string{
			"admin": globalAdminCert,
			"user":  globalUserCert,
		},
		ioutil.ReadFile,
	)

	if err != nil {
		panic(err)
	}

	localActors, err := testcc.IdentitiesFromFiles(
		localMSP,
		map[string]string{
			"admin": localAdminCert,
		},
		ioutil.ReadFile,
	)

	if err != nil {
		panic(err)
	}

	globalUserIdentity := SerializeIdentity(globalMSP, globalActors["user"].Certificate)

	BeforeSuite(func() {
		// Init chaincode before running any tests
		expectcc.ResponseOk(cc.From(localActors["admin"]).Init())
	})

	Describe("Charge", func() {
		It("Allow a global user to start charge transaction", func() {
			startTransaction := &StartTransaction{
				Contractor: contractorName,
			}

			resp := cc.From(globalActors["user"]).Invoke("InvokeStartChargeTransaction", startTransaction)
			invokeResponse := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			queryResp := cc.From(globalActors["user"]).Invoke("QueryById", contractorName, invokeResponse.ChargeId)

			createdTransaction := expectcc.PayloadIs(queryResp, &Entity{}).(Entity)

			Expect(createdTransaction.Contractor).To(Equal(contractorName))
			Expect(createdTransaction.ChargeId).To(Equal(invokeResponse.ChargeId))
			Expect(createdTransaction.State).To(Equal(ChargeStateStarted))
			Expect(createdTransaction.User).To(Equal(globalUserIdentity))
			Expect(createdTransaction.Price).To(BeZero())
			Expect(createdTransaction.EndTime.Second()).To(BeZero())
		})

		It("Allow a global user get a list of its charge transactions", func() {
			// TODO: Limit to only the transactions of the user
			queryResp := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(queryResp, &[]Entity{}).([]Entity)

			Expect(len(chargeTransactions)).To(Equal(1))
		})

		It("Allow a global user to stop charge transaction", func() {
			allUserTransactions := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(allUserTransactions, &[]Entity{}).([]Entity)

			Expect(len(chargeTransactions)).To(Equal(1))

			stopTransaction := &StopTransaction{
				Contractor: contractorName,
				ChargeId:   chargeTransactions[0].ChargeId,
			}

			resp := cc.From(globalActors["user"]).Invoke("InvokeStopTransaction", stopTransaction)
			invokeResponse := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			queryResp := cc.From(globalActors["user"]).Invoke("QueryById", contractorName, invokeResponse.ChargeId)

			stoppedTransaction := expectcc.PayloadIs(queryResp, &Entity{}).(Entity)

			Expect(stoppedTransaction.Contractor).To(Equal(contractorName))
			Expect(stoppedTransaction.ChargeId).To(Equal(chargeTransactions[0].ChargeId))
			Expect(stoppedTransaction.State).To(Equal(ChargeStateStopped))
			Expect(stoppedTransaction.User).To(Equal(globalUserIdentity))
			Expect(stoppedTransaction.Price).To(Equal(0))
			Expect(stoppedTransaction.EndTime.Nanosecond()).To(BeNumerically(">=", stoppedTransaction.StartTime.Nanosecond()))
		})

		It("Allows local admin to get the list of transactions", func() {
			queryResp := cc.From(localActors["admin"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(queryResp, &[]Entity{}).([]Entity)

			Expect(len(chargeTransactions)).To(Equal(1))
		})

		It("Allow a global admin to complete charge transaction", func() {
			allUserTransactions := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(allUserTransactions, &[]Entity{}).([]Entity)

			Expect(len(chargeTransactions)).To(Equal(1))

			startTransaction := &StopTransaction{
				Contractor: contractorName,
				ChargeId:   chargeTransactions[0].ChargeId,
			}

			resp := cc.From(globalActors["admin"]).Invoke("InvokeCompleteTransaction", startTransaction)
			invokeResponse := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			queryResp := cc.From(globalActors["user"]).Invoke("QueryById", contractorName, invokeResponse.ChargeId)

			stoppedTransaction := expectcc.PayloadIs(queryResp, &Entity{}).(Entity)

			Expect(stoppedTransaction.Contractor).To(Equal(contractorName))
			Expect(stoppedTransaction.ChargeId).To(Equal(chargeTransactions[0].ChargeId))
			Expect(stoppedTransaction.State).To(Equal(ChargeStateCompleted))
			Expect(stoppedTransaction.User).To(Equal(globalUserIdentity))
			Expect(stoppedTransaction.Price).To(Equal(0))
			Expect(stoppedTransaction.EndTime.Nanosecond()).To(BeNumerically(">=", stoppedTransaction.StartTime.Nanosecond()))
		})
	})
})
