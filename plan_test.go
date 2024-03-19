package main

import (
	"context"
	"github.com/UniB-e-e/unibee-go-client"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlanApi(t *testing.T) {
	unibee.ApiKey = OpenapiKey
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	{
		// test new plan
		t.Run("Test New Main Plan", func(t *testing.T) {
			resp, httpRes, err := apiClient.Plan.PlanNewPost(ctx).UnibeeApiMerchantPlanNewReq(unibee.UnibeeApiMerchantPlanNewReq{
				PlanName:      "testPlanByApi",
				Amount:        100,
				Currency:      "USD",
				Type:          nil,
				Description:   nil,
				IntervalCount: unibee.PtrInt32(1),
				IntervalUnit:  unibee.String("Day"),
				Metadata:      nil,
				MetricLimits:  nil,
				AddonIds:      nil,
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
		})
	}
	{
		// test edit plan, plan can not edit after activated
		t.Run("Test Edit Main Plan", func(t *testing.T) {
			resp, httpRes, err := apiClient.Plan.PlanEditPost(ctx).UnibeeApiMerchantPlanEditReq(unibee.UnibeeApiMerchantPlanEditReq{
				PlanName:      "testPlanByApi",
				Amount:        100,
				Currency:      "USD",
				Description:   nil,
				IntervalCount: unibee.PtrInt32(1),
				IntervalUnit:  unibee.String("Day"),
				Metadata:      nil,
				MetricLimits:  nil,
				AddonIds:      nil,
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
		})
	}
}
