package serviceauthtoken_test

import (
	"cf"
	. "cf/commands/serviceauthtoken"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestCreateServiceAuthTokenFailsWithUsage(t *testing.T) {
	authTokenRepo := &testapi.FakeAuthTokenRepo{}
	reqFactory := &testreq.FakeReqFactory{}

	ui := callCreateServiceAuthToken(t, []string{}, reqFactory, authTokenRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callCreateServiceAuthToken(t, []string{"arg1"}, reqFactory, authTokenRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callCreateServiceAuthToken(t, []string{"arg1", "arg2"}, reqFactory, authTokenRepo)
	assert.True(t, ui.FailedWithUsage)

	ui = callCreateServiceAuthToken(t, []string{"arg1", "arg2", "arg3"}, reqFactory, authTokenRepo)
	assert.False(t, ui.FailedWithUsage)
}

func TestCreateServiceAuthTokenRequirements(t *testing.T) {
	authTokenRepo := &testapi.FakeAuthTokenRepo{}
	reqFactory := &testreq.FakeReqFactory{}
	args := []string{"arg1", "arg2", "arg3"}

	reqFactory.LoginSuccess = true
	callCreateServiceAuthToken(t, args, reqFactory, authTokenRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	reqFactory.LoginSuccess = false
	callCreateServiceAuthToken(t, args, reqFactory, authTokenRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestCreateServiceAuthToken(t *testing.T) {
	authTokenRepo := &testapi.FakeAuthTokenRepo{}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
	args := []string{"a label", "a provider", "a value"}

	ui := callCreateServiceAuthToken(t, args, reqFactory, authTokenRepo)
	assert.Contains(t, ui.Outputs[0], "Creating service auth token as")
	assert.Contains(t, ui.Outputs[0], "my-user")
	authToken_Auto := cf.ServiceAuthToken{}
	authToken_Auto.Label = "a label"
	authToken_Auto.Provider = "a provider"
	authToken_Auto.Token = "a value"
	assert.Equal(t, authTokenRepo.CreatedServiceAuthToken, authToken_Auto)

	assert.Contains(t, ui.Outputs[1], "OK")
}

func callCreateServiceAuthToken(t *testing.T, args []string, reqFactory *testreq.FakeReqFactory, authTokenRepo *testapi.FakeAuthTokenRepo) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "my-user",
	})
	assert.NoError(t, err)
	org_Auto := cf.Organization{}
	org_Auto.Name = "my-org"
	space_Auto := cf.Space{}
	space_Auto.Name = "my-space"
	config := &configuration.Configuration{
		Space:        space_Auto,
		Organization: org_Auto,
		AccessToken:  token,
	}

	cmd := NewCreateServiceAuthToken(ui, config, authTokenRepo)
	ctxt := testcmd.NewContext("create-service-auth-token", args)

	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
