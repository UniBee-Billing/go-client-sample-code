package main

import (
	"crypto/tls"
	"fmt"
	"github.com/UniB-e-e/unibee-go-client"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	{
		// Websocket use to mock receive webhook message from UniBee
		client := gclient.NewWebSocket()
		client.HandshakeTimeout = time.Second
		client.Proxy = http.ProxyFromEnvironment
		client.TLSClientConfig = &tls.Config{}

		for {
			conn, _, err := client.Dial("ws://api.unibee.top/merchant_ws/"+OpenapiKey, nil)
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
						var stat *unibee.UnibeeInternalLogicGatewayRoUserMetric
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
