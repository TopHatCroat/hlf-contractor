package chaincode

import (
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
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
	cryptoPath   = path.Join("testdata")
	org1UserPub  = path.Join(cryptoPath, "org1", "user", "user.pem")
	org1UserPriv = path.Join(cryptoPath, "org1", "user", "user.key.pem")

	org2UserPub  = path.Join(cryptoPath, "org2", "user", "user.pem")
	org2UserPriv = path.Join(cryptoPath, "org2", "user", "user.key.pem")

	adminPub  = path.Join(cryptoPath, "org1", "admin", "admin.pem")
	adminPriv = path.Join(cryptoPath, "org1", "admin", "admin.key.pem")
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Suite")
}

var _ = Describe("Project", func() {
	r := router.New("root")
	debug.AddHandlers(r, "debug", owner.Only)
	CreateRouter(r)
	r.Init(owner.InvokeSetFromCreator)
	cc := testcc.NewMockStub("project", router.NewChaincode(r))

	// Load actor certificates
	actors, err := testcc.IdentitiesFromFiles(
		"SOME_MSP",
		map[string]string{
			"admin":    adminPub,
			"org1User": org1UserPub,
			"org2User": org2UserPub,
		},
		ioutil.ReadFile,
	)

	if err != nil {
		panic(err)
	}

	org1UserIdentitySerialized := actors["org1User"].Certificate.Subject.CommonName
	org2UserIdentitySerialized := actors["org2User"].Certificate.Subject.CommonName

	BeforeSuite(func() {
		// Init chaincode before running any tests
		expectcc.ResponseOk(cc.From(actors["admin"]).Init())
	})

	Describe("Project", func() {
		It("Allow a everyone to publish a valid project", func() {
			now := ptypes.TimestampNow()
			newProject := &schema.PublishProject{
				Name:           "proj1",
				Assessor:       org2UserIdentitySerialized,
				StartDate:      now,
				EndDate:        now,
				EstimatedValue: 10000,
				Description:    "Some text",
			}

			resp := cc.From(actors["org1User"]).Invoke("ProjectPublish", newProject)
			expectcc.ResponseOk(resp)

			queryResp := cc.From(actors["org1User"]).Invoke("ProjectGet", &schema.ProjectId{
				Issuer: org1UserIdentitySerialized,
				Name:   "proj1",
			})

			createdProject := expectcc.PayloadIs(queryResp, &schema.Project{}).(*schema.Project)

			Expect(createdProject.Name).To(Equal("proj1"))
			Expect(createdProject.Issuer).To(Equal(org1UserIdentitySerialized))
			Expect(createdProject.Assessor).To(Equal(org2UserIdentitySerialized))
			Expect(createdProject.StartDate.Seconds).To(Equal(now.Seconds))
			Expect(createdProject.EndDate.Seconds).To(Equal(now.Seconds))
			Expect(createdProject.EstimatedValue).To(BeNumerically("==", 10000))
			Expect(createdProject.Description).To(Equal("Some text"))
		})

		It("Allows everyone to get the list of projects", func() {
			queryResp := cc.From(actors["org1User"]).Invoke("ProjectList")

			existingProjects := expectcc.PayloadIs(queryResp, &schema.ProjectList{}).(*schema.ProjectList)

			Expect(len(existingProjects.Items)).To(Equal(1))
		})

		It("Allow a user to publish a valid application to a project", func() {
			projectResp := cc.From(actors["org1User"]).Invoke("ProjectGet", &schema.ProjectId{
				Issuer: org1UserIdentitySerialized,
				Name:   "proj1",
			})

			createdProject := expectcc.PayloadIs(projectResp, &schema.Project{}).(*schema.Project)

			newApplication := &schema.PublishApplication{
				ProjectIssuer: createdProject.Issuer,
				ProjectName:   createdProject.Name,
				Price:         100000,
				Description:   "Some description",
			}

			resp := cc.From(actors["org2User"]).Invoke("ApplicationPublish", newApplication)
			expectcc.ResponseOk(resp)

			applicationResp := cc.From(actors["org1User"]).Invoke("ApplicationGet", &schema.ApplicationId{
				Contractor:    org2UserIdentitySerialized,
				ProjectName:   createdProject.Name,
				ProjectIssuer: createdProject.Issuer,
			})

			createdApplication := expectcc.PayloadIs(applicationResp, &schema.Application{}).(*schema.Application)
			Expect(createdApplication.Contractor).To(Equal(org2UserIdentitySerialized))
			Expect(createdApplication.Description).To(Equal("Some description"))
			Expect(createdApplication.ProjectName).To(Equal(createdProject.Name))
			Expect(createdApplication.ProjectIssuer).To(Equal(createdProject.Issuer))
			Expect(createdApplication.Price).To(BeNumerically("==", 100000))
			Expect(createdApplication.State).To(Equal(schema.Application_APPLIED))
		})
	})
})
