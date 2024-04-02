package main

import (
	"context"
	"github.com/UniB-e-e/unibee-go-client"
	"testing"
)

func TestEmail(t *testing.T) {
	unibee.ApiKey = "EUXAgwv3Vcr1PFWt2SgBumMHXn3ImBqM"
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	t.Run("Test Email Setup", func(t *testing.T) {
		apiClient.Email.EmailGatewaySetupPost(ctx).UnibeeApiMerchantEmailGatewaySetupReq(unibee.UnibeeApiMerchantEmailGatewaySetupReq{
			IsDefault:   unibee.Bool(true),
			Data:        "${YOUR_SEND_GRID_KEY}",
			GatewayName: "sendgrid",
		})
	})
}
