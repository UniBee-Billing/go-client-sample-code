package main

import (
	"context"
	"fmt"
	"github.com/UniB-e-e/unibee-go-client"
	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOneTimePayment(t *testing.T) {
	ctx := context.Background()
	// Query user restriction from metric api
	configuration := unibee.NewConfiguration()
	configuration.AddDefaultHeader("Authorization", "Bearer "+OpenapiKey) // This is your test secret API key.
	configuration.Servers = unibee.ServerConfigurations{unibee.ServerConfiguration{
		URL: UniBeeStageUrl,
	}}
	apiClient := unibee.NewAPIClient(configuration)
	{
		t.Run("Test Payment Gateway", func(t *testing.T) {
			resp, httpRes, err := apiClient.Gateway.GatewayListGet(ctx).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			fmt.Printf("Merchant's Gateway List = %s\n", ToJsonString(resp.Data.Gateways))
		})
	}
	{
		t.Run("Test Payment New", func(t *testing.T) {
			gatewayResp, httpRes, err := apiClient.Gateway.GatewayListGet(ctx).Execute()
			require.Nil(t, err)
			require.NotNil(t, gatewayResp)
			assert.Equal(t, 200, httpRes.StatusCode)
			require.NotEmpty(t, gatewayResp.Data.Gateways)
			require.Greater(t, len(gatewayResp.Data.Gateways), 0)
			// test without items
			resp, httpRes, err := apiClient.Payment.PaymentNewPost(ctx).UnibeeApiMerchantPaymentNewReq(unibee.UnibeeApiMerchantPaymentNewReq{
				TotalAmount:       100,
				Currency:          "usd",
				Email:             "jack.fu@wowow.io",
				ExternalPaymentId: uuid.New().String(),
				ExternalUserId:    "1709272139",
				GatewayId:         *gatewayResp.Data.Gateways[0].GatewayId,
				Items:             nil,
				Metadata:          nil,
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			fmt.Printf("Payment Url is %s\n", *resp.Data.Link)
		})
	}
	//{
	//	t.Run("Test Payment New", func(t *testing.T) {
	//		gatewayResp, httpRes, err := apiClient.Gateway.GatewayListGet(ctx).Execute()
	//		require.Nil(t, err)
	//		require.NotNil(t, gatewayResp)
	//		assert.Equal(t, 200, httpRes.StatusCode)
	//		require.NotEmpty(t, gatewayResp.Data.Gateways)
	//		require.Greater(t, len(gatewayResp.Data.Gateways), 0)
	//		// test without items
	//		resp, httpRes, err := apiClient.Payment.PaymentNewPost(ctx).UnibeeApiMerchantPaymentNewReq(openapiclient.UnibeeApiMerchantPaymentNewReq{
	//			TotalAmount:       200,
	//			Currency:          "usd",
	//			Email:             "jack.fu@wowow.io",
	//			ExternalPaymentId: uuid.New().String(),
	//			ExternalUserId:    "1709272139",
	//			GatewayId:         *gatewayResp.Data.Gateways[0].GatewayId,
	//			LineItems: []openapiclient.UnibeeApiMerchantPaymentItem{{
	//				Amount:                 100,
	//				AmountExcludingTax:     nil,
	//				Currency:               nil,
	//				Description:            nil,
	//				Quantity:               nil,
	//				Tax:                    nil,
	//				TaxScale:               nil,
	//				UnitAmountExcludingTax: nil,
	//			}},
	//			RedirectUrl: nil,
	//			Metadata:    nil,
	//		}).Execute()
	//		require.Nil(t, err)
	//		require.NotNil(t, resp)
	//		assert.Equal(t, 200, httpRes.StatusCode)
	//		fmt.Printf("Payment Url is %s\n", *resp.Data.Link)
	//	})
	//}
}
