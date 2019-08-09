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

var _ = Describe("Users", func() {
	r := router.New("root")
	debug.AddHandlers(r, "debug", owner.Only)
	CreateRouter(r)
	r.Init(owner.InvokeSetFromCreator)
	cc := testcc.NewMockStub("Users", router.NewChaincode(r))

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

	var globalUserUsername = globalActors["user"].Certificate.Subject.CommonName

	BeforeSuite(func() {
		expectcc.ResponseOk(cc.From(globalActors["admin"]).Init())
	})

	Describe("Users", func() {
		It("Allow a global admin to get a full user response", func() {
			resp := cc.From(globalActors["admin"]).Invoke("QueryById", globalMSP, globalUserUsername)
			queryResp := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			Expect(queryResp.Email).To(Equal(globalUserUsername))
			Expect(queryResp.State).To(Equal(UserStateActive))
			Expect(queryResp.Balance).To(BeNumerically("==", 0))
		})

		It("Allow a global admin get a list of all users", func() {
			queryResp := cc.From(globalActors["admin"]).Invoke("QueryAll")

			chargeTransactions := expectcc.PayloadIs(queryResp, &[]Entity{}).([]Entity)
			Expect(len(chargeTransactions)).To(Equal(1))
		})

		It("Do not allow user to get list of users", func() {
			queryResp := cc.From(globalActors["user"]).Invoke("QueryAll")
			Expect(queryResp.Payload).To(BeNil())
		})

		It("Allow a global user to get his own full response", func() {
			resp := cc.From(globalActors["user"]).Invoke("QueryById", globalMSP, globalUserUsername)
			queryResp := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			Expect(queryResp.Email).To(Equal(globalUserUsername))
			Expect(queryResp.State).To(Equal(UserStateActive))
			Expect(queryResp.Balance).To(BeNumerically("==", 0))
		})

		It("Allows global admin to block user", func() {
			resp := cc.From(globalActors["admin"]).Invoke("InvokeBlockUserTransaction", globalMSP, globalUserUsername)
			invokeResponse := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			Expect(invokeResponse.MSP).To(Equal(globalMSP))
			Expect(invokeResponse.Email).To(Equal(globalUserUsername))
			Expect(invokeResponse.State).To(Equal(UserStateBlocked))
		})

		It("Allows global admin to unblock user", func() {
			resp := cc.From(globalActors["admin"]).Invoke("InvokeUnblockUserTransaction", globalMSP, globalUserUsername)
			invokeResponse := expectcc.PayloadIs(resp, &Entity{}).(Entity)

			Expect(invokeResponse.MSP).To(Equal(globalMSP))
			Expect(invokeResponse.Email).To(Equal(globalUserUsername))
			Expect(invokeResponse.State).To(Equal(UserStateActive))
		})
	})
})
