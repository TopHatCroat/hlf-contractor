package main

import (
	"bou.ke/monkey"
	"github.com/TopHatCroat/hlf-contractor/chaincode/charger/charge"
	"github.com/TopHatCroat/hlf-contractor/chaincode/charger/service"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/tkuchiki/faketime"
	"io/ioutil"
	"path"
	"testing"
	"time"

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

	chargeDuration      = 5 * time.Minute
	chargeStartTime     = time.Now()
	chargeEndTime       = time.Now().Add(chargeDuration)
	fakeChargeStartTime = faketime.NewFaketimeWithTime(chargeStartTime)
	fakeChargeStopTime  = faketime.NewFaketimeWithTime(chargeEndTime)
	priceCentPerMinute  = 3
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Charge Suite")
}

var _ = Describe("Charge", func() {
	r := router.New("root")
	debug.AddHandlers(r, "debug", owner.Only)
	CreateRouter(r)
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
		expectcc.ResponseOk(cc.From(localActors["admin"]).Init("3"))
	})

	Describe("Charge", func() {
		It("Empty list result returns an empty array", func() {
			queryResp := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(queryResp, &[]charge.Entity{}).([]charge.Entity)

			Expect(chargeTransactions).ToNot(BeNil())
			Expect(len(chargeTransactions)).To(Equal(0))
		})

		It("Throws error when global admin tries to get non existent charge", func() {
			resp := cc.From(globalActors["admin"]).Invoke("QueryById", globalMSP, "doesNotExists")

			Expect(resp.Status).To(BeNumerically("==", 500))
			Expect(resp.Message).To(ContainSubstring("doesNotExists"))
		})

		It("Do not allow a global user to start charge transaction if blocked", func() {
			fakeChargeStartTime.Do()
			guard := mockGetUser("blocked")

			startTransaction := &charge.StartTransaction{
				Contractor: contractorName,
			}

			resp := cc.From(globalActors["user"]).Invoke("InvokeStartChargeTransaction", startTransaction)

			Expect(resp.Status).To(BeNumerically("==", 500))
			Expect(resp.Message).To(ContainSubstring("is blocked"))

			fakeChargeStartTime.Undo()
			guard.Unpatch()
		})

		It("Allow a global user to start charge transaction", func() {
			fakeChargeStartTime.Do()
			guard := mockGetUser("active")

			startTransaction := &charge.StartTransaction{
				Contractor: contractorName,
			}

			resp := cc.From(globalActors["user"]).Invoke("InvokeStartChargeTransaction", startTransaction)
			invokeResponse := expectcc.PayloadIs(resp, &charge.Entity{}).(charge.Entity)

			queryResp := cc.From(globalActors["user"]).Invoke("QueryById", contractorName, invokeResponse.ChargeId)

			createdTransaction := expectcc.PayloadIs(queryResp, &charge.Entity{}).(charge.Entity)

			Expect(createdTransaction.Contractor).To(Equal(contractorName))
			Expect(createdTransaction.ChargeId).To(Equal(invokeResponse.ChargeId))
			Expect(createdTransaction.State).To(Equal(charge.ChargeStateStarted))
			Expect(createdTransaction.User).To(Equal(globalUserIdentity))
			Expect(createdTransaction.Price).To(BeZero())
			Expect(createdTransaction.StartTime.Second()).To(Equal(chargeStartTime.Second()))
			Expect(createdTransaction.EndTime.Second()).To(BeZero())
			fakeChargeStartTime.Undo()
			guard.Unpatch()
		})

		It("Allow a global user get a list of its charge transactions", func() {
			// TODO: Limit to only the transactions of the user
			queryResp := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(queryResp, &[]charge.Entity{}).([]charge.Entity)

			Expect(len(chargeTransactions)).To(Equal(1))
		})

		It("Allow a global user to stop charge transaction", func() {
			fakeChargeStopTime.Do()
			allUserTransactions := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(allUserTransactions, &[]charge.Entity{}).([]charge.Entity)

			Expect(len(chargeTransactions)).To(Equal(1))

			stopTransaction := &charge.StopTransaction{
				Contractor: contractorName,
				ChargeId:   chargeTransactions[0].ChargeId,
			}

			resp := cc.From(globalActors["user"]).Invoke("InvokeStopChargeTransaction", stopTransaction)
			invokeResponse := expectcc.PayloadIs(resp, &charge.Entity{}).(charge.Entity)

			queryResp := cc.From(globalActors["user"]).Invoke("QueryById", contractorName, invokeResponse.ChargeId)

			stoppedTransaction := expectcc.PayloadIs(queryResp, &charge.Entity{}).(charge.Entity)

			Expect(stoppedTransaction.Contractor).To(Equal(contractorName))
			Expect(stoppedTransaction.ChargeId).To(Equal(chargeTransactions[0].ChargeId))
			Expect(stoppedTransaction.State).To(Equal(charge.ChargeStateStopped))
			Expect(stoppedTransaction.User).To(Equal(globalUserIdentity))
			Expect(stoppedTransaction.Price).To(Equal(int(chargeDuration.Minutes()) * priceCentPerMinute))
			Expect(stoppedTransaction.EndTime.Second()).To(Equal(chargeEndTime.Second()))

			fakeChargeStopTime.Undo()
		})

		It("Allows local admin to get the list of transactions", func() {
			queryResp := cc.From(localActors["admin"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(queryResp, &[]charge.Entity{}).([]charge.Entity)

			Expect(len(chargeTransactions)).To(Equal(1))
		})

		It("Allow a global admin to complete charge transaction", func() {
			allUserTransactions := cc.From(globalActors["user"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(allUserTransactions, &[]charge.Entity{}).([]charge.Entity)

			Expect(len(chargeTransactions)).To(Equal(1))

			startTransaction := &charge.StopTransaction{
				Contractor: contractorName,
				ChargeId:   chargeTransactions[0].ChargeId,
			}

			resp := cc.From(globalActors["admin"]).Invoke("InvokeCompleteChargeTransaction", startTransaction)
			invokeResponse := expectcc.PayloadIs(resp, &charge.Entity{}).(charge.Entity)

			queryResp := cc.From(globalActors["user"]).Invoke("QueryById", contractorName, invokeResponse.ChargeId)

			stoppedTransaction := expectcc.PayloadIs(queryResp, &charge.Entity{}).(charge.Entity)

			Expect(stoppedTransaction.Contractor).To(Equal(contractorName))
			Expect(stoppedTransaction.ChargeId).To(Equal(chargeTransactions[0].ChargeId))
			Expect(stoppedTransaction.State).To(Equal(charge.ChargeStateCompleted))
			Expect(stoppedTransaction.User).To(Equal(globalUserIdentity))
			Expect(stoppedTransaction.Price).To(Equal(int(chargeDuration.Minutes()) * priceCentPerMinute))
			Expect(stoppedTransaction.EndTime.Second()).To(Equal(chargeEndTime.Second()))
		})
	})
})

func mockGetUser(state string) *monkey.PatchGuard {
	return monkey.Patch(service.GetUser, func(shim.ChaincodeStubInterface, string, string) (*service.User, error) {
		return &service.User{
			MSP:     globalMSP,
			Email:   "notimportant@mail.com",
			State:   state,
			Balance: 0,
		}, nil
	})
}
