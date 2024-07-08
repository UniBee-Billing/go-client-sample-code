package main

import (
	"context"
	"github.com/UniBee-Billing/unibee-go-client"
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
	var planId int64
	{
		// test new plan
		t.Run("Test New Main Plan", func(t *testing.T) {
			resp, httpRes, err := apiClient.Plan.PlanNewPost(ctx).UnibeeApiMerchantPlanNewReq(unibee.UnibeeApiMerchantPlanNewReq{
				PlanName:      "testPlanByApi",
				Amount:        100,
				Currency:      "USD",
				Type:          unibee.Int32(unibee.PlanTypeMain),
				Description:   nil,
				IntervalCount: unibee.PtrInt32(1),
				IntervalUnit:  unibee.String("Day"),
				Metadata:      &map[string]string{"testKey": "testValue"},
				MetricLimits:  nil,
				AddonIds:      nil,
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			planId = *resp.Data.Plan.Id
		})
	}
	{
		// test edit plan, plan can not edit after activated
		t.Run("Test Edit Main Plan", func(t *testing.T) {
			resp, httpRes, err := apiClient.Plan.PlanEditPost(ctx).UnibeeApiMerchantPlanEditReq(unibee.UnibeeApiMerchantPlanEditReq{
				PlanId:        planId,
				PlanName:      unibee.String("testPlanByApi"),
				Amount:        unibee.Int32(100),
				Currency:      unibee.String("USD"),
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
		t.Run("Test Get Plan List", func(t *testing.T) {
			resp, httpRes, err := apiClient.Plan.PlanListGet(ctx).Page(0).Count(20).Execute()
			if err != nil {
				return
			}
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data.Plans)
			assert.Equal(t, 200, httpRes.StatusCode)
		})
		t.Run("Test Get Plan List Use Post", func(t *testing.T) {
			resp, httpRes, err := apiClient.Plan.PlanListPost(ctx).UnibeeApiMerchantPlanListReq(unibee.UnibeeApiMerchantPlanListReq{
				Count:         unibee.PtrInt32(20),
				Page:          unibee.PtrInt32(0),
				Currency:      nil,
				PublishStatus: nil,
				SortField:     nil,
				SortType:      nil,
				Status:        nil,
				Type:          nil,
			}).Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data.Plans)
			assert.Equal(t, 200, httpRes.StatusCode)
		})
	}
	{
		apiClient.Plan.PlanActivatePost(ctx).UnibeeApiMerchantPlanActivateReq(unibee.UnibeeApiMerchantPlanActivateReq{PlanId: planId})
	}
}
