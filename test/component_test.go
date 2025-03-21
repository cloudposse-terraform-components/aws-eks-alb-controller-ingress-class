package test

import (
	"context"
	"testing"
	"fmt"
	"strings"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	awsHelper "github.com/cloudposse/test-helpers/pkg/aws"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	// "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/gruntwork-io/terratest/modules/random"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "eks/alb-controller-ingress-class/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	randomID := strings.ToLower(random.UniqueId())
	class_name := fmt.Sprintf("alb-%s", randomID)
	group_name := fmt.Sprintf("group-%s", randomID)

	inputs := map[string]interface{}{
		"class_name": class_name,
		"group": group_name,
		"ip_address_type": "ipv4",
		"scheme": "internet-facing",
	}

	defer s.DestroyAtmosComponent(s.T(), component, stack, &inputs)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, &inputs)
	assert.NotNil(s.T(), options)

	clusterOptions := s.GetAtmosOptions("eks/cluster", stack, nil)
	clusrerId := atmos.Output(s.T(), clusterOptions, "eks_cluster_id")
	cluster := awsHelper.GetEksCluster(s.T(), context.Background(), awsRegion, clusrerId)
	clientset, err := awsHelper.NewK8SClientset(cluster)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), clientset)

	ingressClass, err := clientset.NetworkingV1().IngressClasses().Get(context.Background(), class_name, metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ingressClass.ObjectMeta.Name, class_name)
	assert.Equal(s.T(), ingressClass.Spec.Controller, "ingress.k8s.aws/alb")

	s.DriftTest(component, stack, &inputs)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "eks/alb-controller-ingress-class/disabled"
	const stack = "default-test"
	s.VerifyEnabledFlag(component, stack, nil)
}

func (s *ComponentSuite) SetupSuite() {
	s.TestSuite.InitConfig()
	s.TestSuite.Config.ComponentDestDir = "components/terraform/eks/alb-controller-ingress-class"
	s.TestSuite.SetupSuite()
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)
	suite.AddDependency(t, "vpc", "default-test", nil)
	suite.AddDependency(t, "eks/cluster", "default-test", nil)
	suite.AddDependency(t, "eks/alb-controller", "default-test", nil)
	helper.Run(t, suite)
}
