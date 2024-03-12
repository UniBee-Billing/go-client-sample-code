package main

import (
	"context"
	"fmt"
	"github.com/UniB-e-e/unibee-go-client"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetricApi(t *testing.T) {
	unibee.ApiKey = OpenapiKey
	unibee.Host = UniBeeStageUrl
	ctx := context.Background()
	// Query user restriction from metric api
	configuration := unibee.NewConfiguration()
	apiClient := unibee.NewAPIClient(configuration)
	{
		t.Run("Test Metric GetUserMetric", func(t *testing.T) {
			resp, httpRes, err := apiClient.UserMetric.MetricUserMetricGet(ctx).ExternalUserId("1709272139").Execute()
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)
			fmt.Printf("Unibee's User isPaid = %v\n", *resp.Data.GetUserMetric().IsPaid)
			if resp.Data.GetUserMetric().Plan != nil {
				fmt.Printf("Unibee's User PlanName = %s\n", *resp.Data.GetUserMetric().Plan.PlanName)
				extraData := gjson.New(*resp.Data.GetUserMetric().Plan.ExtraMetricData)
				fmt.Printf("allowed_browser_types %s\n", extraData.Get("allowed_browser_types"))
			}
			var userRestrictionMap = make(map[string]interface{})
			for _, metric := range resp.Data.GetUserMetric().UserMerchantMetricStats {
				userRestrictionMap[*metric.MetricLimit.Code] = *metric.MetricLimit.TotalLimit
			}
			fmt.Printf("Unibee's Metric Limit List = %v\n", userRestrictionMap)
		})
	}
}
