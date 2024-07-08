package main

import (
	"context"
	"github.com/UniBee-Billing/unibee-go-client"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVatSetup(t *testing.T) {
	unibee.ApiKey = "5fshy9X6e3Wn0a6xjyRLMq2GGNM4sTEf"
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	t.Run("Test Vat Setup", func(t *testing.T) {
		resp, httpRes, err := apiClient.Vat.VatSetupGatewayPost(ctx).UnibeeApiMerchantVatSetupGatewayReq(unibee.UnibeeApiMerchantVatSetupGatewayReq{
			IsDefault:   unibee.Bool(true),
			Data:        "d9d57f2212cba7e286b3fb9cbb2ad419",
			GatewayName: "vatsense",
		}).Execute()
		require.Nil(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 200, httpRes.StatusCode)
	})
}
