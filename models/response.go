package models

import "time"

// APIResponse 定义统一的API响应格式
type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

// ResponseSuccess 返回成功响应
func ResponseSuccess(data interface{}) APIResponse {
	return APIResponse{
		Code: 200,
		Msg:  "success",
		Data: data,
		Time: time.Now().Unix(),
	}
}

// ResponseError 返回错误响应
func ResponseError(code int, msg string) APIResponse {
	return APIResponse{
		Code: code,
		Msg:  msg,
		Data: nil,
		Time: time.Now().Unix(),
	}
}
