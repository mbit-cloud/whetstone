package whetstone_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry-incubator/runtime-schema/models/factories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego Edge", func() {
	Context("when desiring a docker-based LRP", func() {
		It("eventually runs on an executor", func() {
			processGuid := factories.GenerateGuid()

			repUrl := "http://127.0.0.1:20515"

			err := bbs.DesireLRP(models.DesiredLRP{
				Domain:      "whettest-stone-ever",
				ProcessGuid: processGuid,
				Instances:   1,
				Stack:       "lucid64",
				RootFSPath:  "docker:///onsi/grace-busybox",
				Routes: []string{"whetgrace.10.244.0.34.xip.io"},
				MemoryMB: 128,
				DiskMB:   1024,
				Ports: []models.PortMapping{
					{ContainerPort: 8080},
				},
				Actions: []models.ExecutorAction{
					models.Parallel(
						models.ExecutorAction{
							models.RunAction{
								Path: "/app",
								//								app expects $VCAP_APPLICATION
								Env: []models.EnvironmentVariable{
									{Name: "VCAP_APPLICATION", Value: `{"instance_index":0}`},
									{Name: "PORT", Value: "8080"},
								},
							},
						},
						models.ExecutorAction{
							models.MonitorAction{
								Action: models.ExecutorAction{
									models.RunAction{ //The spy. Is this container healthy? running on 8080?
										Path: "echo",
										Args: []string{"all good"},
									},
								},
								HealthyThreshold:   1,
								UnhealthyThreshold: 1,
								HealthyHook: models.HealthRequest{ //Teel the rep where to call back to on exit 0 of spy
									Method: "PUT",
									URL: fmt.Sprintf(
										"http://%s/lrp_running/%s/PLACEHOLDER_INSTANCE_INDEX/PLACEHOLDER_INSTANCE_GUID",
										repUrl,
										processGuid,
									),
								},
							},
						},
					),
				},
			})

			Ω(err).ShouldNot(HaveOccurred())

		})
	})
})
