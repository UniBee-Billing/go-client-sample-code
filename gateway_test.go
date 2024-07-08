package main

import (
	"context"
	"fmt"
	"github.com/UniBee-Billing/unibee-go-client"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetGatewayWebhookUrl(t *testing.T) {
	unibee.ApiKey = "EUXAgwv3Vcr1PFWt2SgBumMHXn3ImBqM"
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	{
		t.Run("Test Changelly GetGatewayWebhookUrl", func(t *testing.T) {
			resp, httpRes, err := apiClient.Gateway.GatewaySetupWebhookPost(ctx).UnibeeApiMerchantGatewaySetupWebhookReq(unibee.UnibeeApiMerchantGatewaySetupWebhookReq{GatewayId: 29}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			fmt.Printf("GatewayWebhookUrl %s\n", *resp.Data.GatewayWebhookUrl)
		})
		t.Run("Test Changelly Setup", func(t *testing.T) {
			apiClient.Gateway.GatewaySetupPost(ctx).UnibeeApiMerchantGatewaySetupReq(unibee.UnibeeApiMerchantGatewaySetupReq{
				GatewayName:   "",
				GatewayKey:    nil,
				GatewaySecret: nil,
			})
		})
	}
}
