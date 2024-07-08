package main

import (
	"context"
	"fmt"
	"github.com/UniBee-Billing/unibee-go-client"
	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOneTimePayment(t *testing.T) {
	unibee.ApiKey = OpenapiKey
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	configuration := unibee.NewConfiguration()
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
			// get the gateway, gateway need setup int merchant portal or use api
			gatewayResp, httpRes, err := apiClient.Gateway.GatewayListGet(ctx).Execute()
			require.Nil(t, err)
			require.NotNil(t, gatewayResp)
			assert.Equal(t, 200, httpRes.StatusCode)
			require.NotEmpty(t, gatewayResp.Data.Gateways)
			require.Greater(t, len(gatewayResp.Data.Gateways), 0)
			// create new payment
			resp, httpRes, err := apiClient.Payment.PaymentNewPost(ctx).UnibeeApiMerchantPaymentNewReq(unibee.UnibeeApiMerchantPaymentNewReq{
				CountryCode:       unibee.String("CH"),
				TotalAmount:       unibee.Int64(100),
				Currency:          unibee.String("USD"),
				Email:             unibee.String("jack.fu@wowow.io"),
				ExternalPaymentId: unibee.String(uuid.New().String()), // your paymentId
				ExternalUserId:    unibee.String("1709272139"),
				GatewayId:         *gatewayResp.Data.Gateways[0].GatewayId,
				Items:             nil, // without items
				Metadata:          &map[string]string{"key1": "value1", "key2": "value2"},
				RedirectUrl:       unibee.String("http://localhost/paymentResult"),
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			assert.Equal(t, unibee.PaymentCreated, int(*resp.Data.Status))
			fmt.Printf("Payment Url is %s\n", *resp.Data.Link)
			// get the payment detail
			detailResp, detailHttpRes, err := apiClient.Payment.PaymentDetailGet(ctx).PaymentId(unibee.StringValue(resp.Data.PaymentId)).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, detailHttpRes.StatusCode)
			require.NotNil(t, detailResp.Data.PaymentDetail)
			assert.Equal(t, unibee.PaymentCreated, int(*detailResp.Data.PaymentDetail.Payment.Status))
			fmt.Printf("Payment Detail is %s", ToJsonString(detailResp.Data.PaymentDetail))
		})
	}
	{
		t.Run("Test Payment Refund", func(t *testing.T) {
			// refund a payment, refund can call more than one util payment refunded completely
			resp, httpRes, err := apiClient.Payment.PaymentRefundNewPost(ctx).UnibeeApiMerchantPaymentNewPaymentRefundReq(unibee.UnibeeApiMerchantPaymentNewPaymentRefundReq{
				PaymentId:        "pay20240310d2Rv1eXrQq2a6jA",
				ExternalRefundId: uuid.New().String(), // your refundId
				RefundAmount:     100,
				Currency:         "USD",
				Metadata:         &map[string]string{"key1": "value1", "key2": "value2"},
				Reason:           unibee.String("refund test"),
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			detailResp, detailHttpRes, err := apiClient.Payment.PaymentRefundDetailGet(ctx).RefundId(unibee.StringValue(resp.Data.RefundId)).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, detailHttpRes.StatusCode)
			require.NotNil(t, detailResp.Data.RefundDetail)
			assert.Equal(t, unibee.RefundCreated, int(*resp.Data.Status))
			fmt.Printf("Payment Refund Detail is %s", ToJsonString(detailResp.Data.RefundDetail))
		})
	}
	{
		t.Run("Test Payment New And Cancel", func(t *testing.T) {
			// get the gateway, gateway need setup int merchant portal or use api
			gatewayResp, httpRes, err := apiClient.Gateway.GatewayListGet(ctx).Execute()
			require.Nil(t, err)
			require.NotNil(t, gatewayResp)
			assert.Equal(t, 200, httpRes.StatusCode)
			require.NotEmpty(t, gatewayResp.Data.Gateways)
			require.Greater(t, len(gatewayResp.Data.Gateways), 0)
			// create new payment
			resp, httpRes, err := apiClient.Payment.PaymentNewPost(ctx).UnibeeApiMerchantPaymentNewReq(unibee.UnibeeApiMerchantPaymentNewReq{
				TotalAmount:       unibee.Int64(200),
				Currency:          unibee.String("USD"),
				Email:             unibee.String("jack.fu@wowow.io"),
				ExternalPaymentId: unibee.String(uuid.New().String()),
				ExternalUserId:    unibee.String("1709272139"),
				GatewayId:         *gatewayResp.Data.Gateways[0].GatewayId,
				Items: []unibee.UnibeeApiMerchantPaymentItem{{
					Amount:                 200,
					Description:            unibee.String("test item"),
					Currency:               unibee.String("usd"),
					AmountExcludingTax:     nil,
					Quantity:               nil,
					Tax:                    nil,
					TaxPercentage:          nil,
					UnitAmountExcludingTax: nil,
				}}, // with items
				RedirectUrl: unibee.String("http://localhost/paymentResult"),
				Metadata:    nil,
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			require.NotNil(t, resp.Data.PaymentId)
			require.NotNil(t, resp.Data.Link)
			require.Equal(t, unibee.PaymentCreated, int(*resp.Data.Status))
			fmt.Printf("Payment Url is %s\n", *resp.Data.Link)
			//cancel the payment
			_, cancelHttpRes, err := apiClient.Payment.PaymentCancelPost(ctx).UnibeeApiMerchantPaymentCancelReq(unibee.UnibeeApiMerchantPaymentCancelReq{
				PaymentId: unibee.StringValue(resp.Data.PaymentId),
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, cancelHttpRes.StatusCode)
			detailResp, detailHttpRes, err := apiClient.Payment.PaymentDetailGet(ctx).PaymentId(unibee.StringValue(resp.Data.PaymentId)).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, detailHttpRes.StatusCode)
			require.NotNil(t, detailResp.Data.PaymentDetail)
			assert.Equal(t, unibee.PaymentCancelled, int(*detailResp.Data.PaymentDetail.Payment.Status))
		})
	}

}
