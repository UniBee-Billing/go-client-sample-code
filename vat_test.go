package main

import (
	"context"
	"github.com/UniB-e-e/unibee-go-client"
	"testing"
)

func TestVatSetup(t *testing.T) {
	unibee.ApiKey = "EUXAgwv3Vcr1PFWt2SgBumMHXn3ImBqM"
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	t.Run("Test Vat Setup", func(t *testing.T) {
		apiClient.Vat.VatSetupGatewayPost(ctx).UnibeeApiMerchantVatSetupGatewayReq(unibee.UnibeeApiMerchantVatSetupGatewayReq{
			IsDefault:   unibee.Bool(true),
			Data:        "${YOUR_VAT_SENSE_KEY}",
			GatewayName: "vatsense",
		})
	})
}
