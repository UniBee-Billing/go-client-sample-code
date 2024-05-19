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
			conn, _, err := client.Dial("wss://api.unibee.top/merchant_ws/"+OpenapiKey, nil)
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
					if one != nil {
						fmt.Printf("Receive Webhook event = %s\n", one.WebhookEvent)
						if one.WebhookEvent == "user.metric.update" {
							var metric *unibee.UnibeeApiMerchantMetricUserMetric
							err = gjson.Unmarshal([]byte(one.Data), &metric)
							if err != nil {
								fmt.Printf("error:%s\n", err.Error())
								continue
							}
							fmt.Printf("Unibee's User Metric = %s\n", ToJsonString(metric))
							fmt.Printf("Unibee's User isPaid = %v\n", unibee.BoolValue(metric.IsPaid))
							if metric.Plan != nil {
								fmt.Printf("Unibee's User PlanName = %s\n", unibee.StringValue(metric.Plan.PlanName))
							}
							var userRestrictionMap = make(map[string]interface{})
							for _, metric := range metric.UserMerchantMetricStats {
								userRestrictionMap[unibee.StringValue(metric.MetricLimit.Code)] = unibee.Int64Value(metric.MetricLimit.TotalLimit)
							}
							fmt.Printf("Unibee's Metric Limit List = %v\n", userRestrictionMap)
							// next setup your user's restriction logic
						} else if one.WebhookEvent == "payment.created" ||
							one.WebhookEvent == "payment.success" ||
							one.WebhookEvent == "payment.cancelled" ||
							one.WebhookEvent == "payment.failure" ||
							one.WebhookEvent == "payment.authorised.need" {
							var paymentDetail *unibee.UnibeeApiBeanDetailPaymentDetail
							err = gjson.Unmarshal([]byte(one.Data), &paymentDetail)
							if err != nil {
								fmt.Printf("error:%s\n", err.Error())
								continue
							}
							fmt.Printf("Unibee's Payment = %s\n", ToJsonString(paymentDetail))
						} else if one.WebhookEvent == "refund.created" ||
							one.WebhookEvent == "refund.success" ||
							one.WebhookEvent == "refund.failure" ||
							one.WebhookEvent == "refund.reverse" {
							var refundDetail *unibee.UnibeeApiBeanDetailRefundDetail
							err = gjson.Unmarshal([]byte(one.Data), &refundDetail)
							if err != nil {
								fmt.Printf("error:%s\n", err.Error())
								continue
							}
							fmt.Printf("Unibee's Refund = %s\n", ToJsonString(refundDetail))
						}
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
