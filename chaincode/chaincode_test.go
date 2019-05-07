package chaincode

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/golang/protobuf/ptypes"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/TopHatCroat/hlf-contractor/chaincode/schema"
)

var (
	cryptoPath = path.Join("..", "network", "crypto-config", "peerOrganizations", "org1.example.com")
	userPub    = path.Join(cryptoPath, "users", "User1@org1.example.com", "msp", "admincerts", "User1@org1.example.com-cert.pem")
	userPriv   = path.Join(cryptoPath, "users", "User1@org1.example.com", "msp", "keystore", "c75bd6911aca808941c3557ee7c97e90f3952e379497dc55eb903f31b50abc83_sk")
	adminPub   = path.Join(cryptoPath, "users", "Admin@org1.example.com", "msp", "admincerts", "Admin@org1.example.com-cert.pem")
	adminPriv  = path.Join(cryptoPath, "users", "Admin@org1.example.com", "msp", "keystore", "cd96d5260ad4757551ed4a5a991e62130f8008a0bf996e4e4b84cd097a747fec_sk")
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Suite")
}

var _ = Describe(`Project`, func() {
	cc := testcc.NewMockStub(`project`, NewCC())

	// Load actor certificates
	actors, err := testcc.IdentitiesFromFiles(`SOME_MSP`, map[string]string{
		`authority`: adminPub,
		`someone`:   userPub},
		ioutil.ReadFile,
	)
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		// Init chaincode before running any tests
		expectcc.ResponseOk(cc.From(actors[`authority`]).Init())
	})

	Describe("Project", func() {

		It("Allow a user to publish a valid project", func() {
			now := ptypes.TimestampNow()
			newProject := &schema.PublishProject{
				Name:           "proj1",
				Assessor:       "user@example.com",
				StartDate:      now,
				EndDate:        now,
				EstimatedValue: 10000,
				Description:    "Some text",
			}

			resp := cc.From(actors["someone"]).Invoke("projectPublish", newProject)
			expectcc.ResponseOk(resp)

			queryResp := cc.From(actors["someone"]).Invoke("projectGet", &schema.ProjectId{
				Issuer: "User1@org1.example.com",
				Name: "proj1",
			})

			createdProject := expectcc.PayloadIs(queryResp, &schema.Project{}).(*schema.Project)

			Expect(createdProject.Name).To(Equal("proj1"))
			Expect(createdProject.Assessor).To(Equal("user@example.com"))
			Expect(createdProject.StartDate.Seconds).To(Equal(now.Seconds))
			Expect(createdProject.EndDate.Seconds).To(Equal(now.Seconds))
			Expect(createdProject.EstimatedValue).To(BeNumerically("==", 10000))
			Expect(createdProject.Description).To(Equal("Some text"))
		})
	})
})
