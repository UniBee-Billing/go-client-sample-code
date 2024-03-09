package main

import (
	"context"
	"crypto/tls"
	"fmt"
	openapiclient "github.com/UniB-e-e/unibee-go-client"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func main() {
	openapiKey := "FiJG2DxPE63G2Chy92uZ42TIbokUpaZG" // MLX's openapiKey (stage)
	unibeeStageUrl := "http://api.unibee.top"        // unibee api host (stage)
	env := "stage"

	ctx := context.Background()
	{
		// Create an account in billing and Provides one time URL to authenticate in a billing client area
		configuration := openapiclient.NewConfiguration()
		configuration.AddDefaultHeader("Authorization", "Bearer "+openapiKey) // This is your test secret API key.
		configuration.Servers = openapiclient.ServerConfigurations{openapiclient.ServerConfiguration{
			URL:         unibeeStageUrl,
			Description: env,
		}}
		apiClient := openapiclient.NewAPIClient(configuration)
		resp, _, err := apiClient.MerchantSessionAPI.MerchantSessionNewSessionPost(ctx).UnibeeApiMerchantSessionNewReq(openapiclient.UnibeeApiMerchantSessionNewReq{
			Email:          "jack.fu@wowow.io", // should change to MLS‘s user email
			ExternalUserId: "1709272139",       // should change to ID of MLS‘s user
		}).Execute()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Unibee's Userid = %s\n", *resp.Data.UserId)
			fmt.Printf("Client Auth Url = %s\n", *resp.Data.Url)
		}
	}

	{
		// Query user restriction from metric api
		configuration := openapiclient.NewConfiguration()
		configuration.AddDefaultHeader("Authorization", "Bearer "+openapiKey) // This is your test secret API key.
		configuration.Servers = openapiclient.ServerConfigurations{openapiclient.ServerConfiguration{
			URL:         unibeeStageUrl,
			Description: env,
		}}
		apiClient := openapiclient.NewAPIClient(configuration)
		resp, _, err := apiClient.MerchantUserMetricAPI.MerchantMetricUserMetricGet(ctx).ExternalUserId("1709272139").Execute()
		if err != nil {
			fmt.Println(err)
		} else {
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
		}
	}

	{
		// Websocket receive user restriction from subscription message
		client := gclient.NewWebSocket()
		client.HandshakeTimeout = time.Second
		client.Proxy = http.ProxyFromEnvironment
		client.TLSClientConfig = &tls.Config{}

		for {
			conn, _, err := client.Dial("ws://api.unibee.top/merchant_ws/"+openapiKey, nil)
			if err != nil {
				panic(err)
			}
			for {
				mt, data, err := conn.ReadMessage()
				if err != nil {
					fmt.Printf("error:%s\n", err.Error())
					_ = conn.Close()
					break
				}
				if mt == websocket.BinaryMessage {
					var one *WebhookMessage
					err = gjson.Unmarshal(data, &one)
					if err != nil {
						fmt.Printf("error:%s\n", err.Error())
						continue
					}
					if one != nil && one.WebhookEvent == "user.metric.update" {
						var stat *openapiclient.UnibeeInternalLogicGatewayRoUserMetric
						err = gjson.Unmarshal([]byte(one.Data), &stat)
						if err != nil {
							fmt.Printf("error:%s\n", err.Error())
							continue
						}
						fmt.Printf("Unibee's User isPaid = %v\n", *stat.IsPaid)
						if stat.Plan != nil {
							fmt.Printf("Unibee's User PlanName = %s\n", *stat.Plan.PlanName)
							extraData := gjson.New(*stat.Plan.ExtraMetricData)
							fmt.Printf("allowed_browser_types %s\n", extraData.Get("allowed_browser_types"))
						}
						var userRestrictionMap = make(map[string]interface{})
						for _, metric := range stat.UserMerchantMetricStats {
							userRestrictionMap[*metric.MetricLimit.Code] = *metric.MetricLimit.TotalLimit
						}
						fmt.Printf("Unibee's Metric Limit List = %v\n", userRestrictionMap)
						// next call restriction api
					}
				} else if mt == websocket.PingMessage {
					// ignore ping message
				} else {
					// ignore
				}
			}
		}
	}

}

type WebhookMessage struct {
	Id           uint64 `json:"id"              description:"id"`              // id
	MerchantId   uint64 `json:"merchantId"      description:"merchantId"`      // merchantId
	WebhookEvent string `json:"webhookEvent"    description:"webhook_event"`   // webhook_event
	Data         string `json:"data"            description:"data(json)"`      // data(json)
	CreateTime   int64  `json:"createTime"      description:"create utc time"` // create utc time
}
