package main

import (
	"context"
	"fmt"
	"github.com/UniB-e-e/unibee-go-client"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSessionApi(t *testing.T) {
	unibee.ApiKey = OpenapiKey
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	// Create an account in billing and Provides one time URL to authenticate in a billing client area
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	{
		t.Run("Test Session NewSession", func(t *testing.T) {
			resp, httpRes, err := apiClient.Session.SessionNewSessionPost(ctx).UnibeeApiMerchantSessionNewReq(unibee.UnibeeApiMerchantSessionNewReq{
				Email:          "jack.fu@wowow.io", // should change to MLS‘s user email
				ExternalUserId: "1709272139",       // should change to ID of MLS‘s user
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			require.NotNil(t, resp.Data.Url)
			fmt.Printf("Unibee's Userid = %s\n", *resp.Data.UserId)
			fmt.Printf("Client Auth Url = %s\n", *resp.Data.Url)
		})
	}
}
