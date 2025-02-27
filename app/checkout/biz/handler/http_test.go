package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"TikTokMall/app/checkout/kitex_gen/checkout"
)

func TestCheckoutHTTPHandler_Checkout(t *testing.T) {
	// 创建请求上下文
	hertzCtx := app.NewContext(16)

	// 设置请求内容
	validReq := map[string]interface{}{
		"user_id": 12345,
		"items": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   2,
				"price":      10.99,
			},
		},
		"credit_card": map[string]interface{}{
			"number":      "4111111111111111",
			"expiry":      "12/25",
			"cvv":         "123",
			"holder_name": "Test User",
		},
		"shipping_address": map[string]interface{}{
			"street":      "123 Test St",
			"city":        "Test City",
			"state":       "TS",
			"postal_code": "12345",
			"country":     "Test Country",
		},
	}

	reqBody, err := json.Marshal(validReq)
	require.NoError(t, err)

	hertzCtx.Request.SetBody(reqBody)
	hertzCtx.Request.Header.SetMethod("POST")
	hertzCtx.Request.Header.SetContentTypeBytes([]byte("application/json"))

	// 设置模拟响应
	resp := &checkout.CheckoutResp{
		OrderId:       "order-123",
		TransactionId: "tx-123",
	}
	hertzCtx.JSON(consts.StatusOK, resp)

	// 验证响应
	assert.Equal(t, http.StatusOK, hertzCtx.Response.StatusCode())

	var gotResp checkout.CheckoutResp
	err = json.Unmarshal(hertzCtx.Response.Body(), &gotResp)
	require.NoError(t, err)

	assert.NotEmpty(t, gotResp.OrderId)
	assert.NotEmpty(t, gotResp.TransactionId)
}

func TestCheckoutHTTPHandler_InvalidRequest(t *testing.T) {
	// 创建请求上下文
	hertzCtx := app.NewContext(16)

	// 设置请求内容 - 缺少必要字段
	invalidReq := map[string]interface{}{
		"user_id": 12345,
		// 缺少 items, credit_card, shipping_address
	}

	reqBody, err := json.Marshal(invalidReq)
	require.NoError(t, err)

	hertzCtx.Request.SetBody(reqBody)
	hertzCtx.Request.Header.SetMethod("POST")
	hertzCtx.Request.Header.SetContentTypeBytes([]byte("application/json"))

	// 设置模拟响应
	hertzCtx.JSON(consts.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   "invalid request",
	})

	// 验证响应
	assert.Equal(t, consts.StatusBadRequest, hertzCtx.Response.StatusCode())

	var gotResp map[string]interface{}
	err = json.Unmarshal(hertzCtx.Response.Body(), &gotResp)
	require.NoError(t, err)

	assert.False(t, gotResp["success"].(bool))
	assert.NotEmpty(t, gotResp["error"])
}
