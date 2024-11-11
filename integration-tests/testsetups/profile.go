package testsetups

//revive:disable:dot-imports
import (
	"time"

	. "github.com/onsi/gomega"
	"golang.org/x/sync/errgroup"

	"github.com/goplugin/plugin-testing-framework/blockchain"
	"github.com/goplugin/plugin-testing-framework/k8s/environment"
	reportModel "github.com/goplugin/plugin-testing-framework/testreporters"

	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/testreporters"
)

// PluginProfileTest runs a piece of code on Plugin nodes with PPROF enabled, then downloads the PPROF results
type PluginProfileTest struct {
	Inputs       PluginProfileTestInputs
	TestReporter testreporters.PluginProfileTestReporter
	env          *environment.Environment
	c            blockchain.EVMClient
}

// PluginProfileTestInputs are the inputs necessary to run a profiling tests
type PluginProfileTestInputs struct {
	ProfileFunction func(*client.PluginClient)
	ProfileDuration time.Duration
	PluginNodes  []*client.PluginK8sClient
}

// NewPluginProfileTest prepares a new keeper Plugin profiling test to be run
func NewPluginProfileTest(inputs PluginProfileTestInputs) *PluginProfileTest {
	return &PluginProfileTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (c *PluginProfileTest) Setup(env *environment.Environment) {
	c.ensureInputValues()
	c.env = env
}

// Run runs the profiling test
func (c *PluginProfileTest) Run() {
	profileGroup := new(errgroup.Group)
	for ni, cl := range c.Inputs.PluginNodes {
		pluginNode := cl
		nodeIndex := ni
		profileGroup.Go(func() error {
			profileResults, err := pluginNode.Profile(c.Inputs.ProfileDuration, c.Inputs.ProfileFunction)
			profileResults.NodeIndex = nodeIndex
			if err != nil {
				return err
			}
			c.TestReporter.Results = append(c.TestReporter.Results, profileResults)
			return nil
		})
	}
	Expect(profileGroup.Wait()).ShouldNot(HaveOccurred(), "Error while gathering plugin Profile tests")
}

// Networks returns the networks that the test is running on
func (c *PluginProfileTest) TearDownVals() (*environment.Environment, []*client.PluginK8sClient, reportModel.TestReporter, blockchain.EVMClient) {
	return c.env, c.Inputs.PluginNodes, &c.TestReporter, c.c
}

// ensureValues ensures that all values needed to run the test are present
func (c *PluginProfileTest) ensureInputValues() {
	Expect(c.Inputs.ProfileFunction).ShouldNot(BeNil(), "Forgot to provide a function to profile")
	Expect(c.Inputs.ProfileDuration.Seconds()).Should(BeNumerically(">=", 1), "Time to profile should be at least 1 second")
	Expect(c.Inputs.PluginNodes).ShouldNot(BeNil(), "Plugin nodes you want to profile should be provided")
	Expect(len(c.Inputs.PluginNodes)).Should(BeNumerically(">", 0), "No Plugin nodes provided to profile")
}
