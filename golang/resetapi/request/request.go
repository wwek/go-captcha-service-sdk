package request

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/wenlng/go-captcha-service-sdk/golang/resetapi/errs"
)

var (
	CONTENT_URLENCODED_TYPE = "application/x-www-form-urlencoded"
	CONTENT_JSON_TYPE       = "application/json"
	CONTENT_FORM_DATA_TYPE  = "multipart/form-data"
	CONTENT_STREAM_TYPE     = "application/octet-stream"
	CONTENT_BINARY_TYPE     = "application/binary"
)

var (
	RESP_RESULT_FIELD_CODE    = "code"
	RESP_RESULT_FIELD_MESSAGE = "message"
	RESP_RESULT_FIELD_DATA    = "data"
)

// BaseRespResult 公共响应数据
type BaseRespResult struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// BaseReqHeaderParams 公共请求头参数
type BaseReqHeaderParams struct {
	XAPIKey string `json:"X-API-Key"`
}

func GetHttpClient() *resty.Client {
	return resty.New()
}

func DefaultRequestTimeout(timeout int64) time.Duration {
	return time.Duration(timeout) * time.Millisecond
}

// ParseResp 解析响应的 data 数据
func ParseResp(resp *resty.Response) (data *BaseRespResult, err error) {
	respData := string(resp.Body())
	var resResult = &BaseRespResult{
		Code:    gjson.Get(respData, RESP_RESULT_FIELD_CODE).Int(),
		Message: gjson.Get(respData, RESP_RESULT_FIELD_MESSAGE).String(),
		Data:    gjson.Get(respData, RESP_RESULT_FIELD_DATA).String(),
	}

	if !errs.CheckBizCodeSuccess(resResult.Code) {
		if err = errs.CheckBizCodeErr(resResult.Code); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("code[%d] message:%v", resResult.Code, resResult.Message)
	}

	return resResult, nil
}

// IsHttpStatusSuccess http status code == 200
func IsHttpStatusSuccess(status int) bool {
	return status == http.StatusOK
}

// IsHttpStatusUnauthorized http status code == 401
func IsHttpStatusUnauthorized(code int) bool {
	return code == http.StatusUnauthorized
}

// IsHttpStatusForbidden http status code == 403
func IsHttpStatusForbidden(code int) bool {
	return code == http.StatusForbidden
}
