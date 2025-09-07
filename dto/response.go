package dto

import (
	"encoding/json"
)

// MsmpResponseSuccess 成功响应结构
type MsmpResponseSuccess struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

// MsmpResponseError 错误响应结构
type MsmpResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// MsmpResponseFailure 失败响应结构
type MsmpResponseFailure struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      int               `json:"id"`
	Error   MsmpResponseError `json:"error"`
}

// MsmpResponse 定义响应接口
type MsmpResponse interface {
	IsSuccess() bool
	GetID() int
	GetJSONRPC() string
}

// IsSuccess 判断响应是否为成功响应
func (resp MsmpResponseSuccess) IsSuccess() bool {
	return true
}

// GetID 获取请求ID
func (resp MsmpResponseSuccess) GetID() int {
	return resp.ID
}

// GetJSONRPC 获取JSON-RPC版本
func (resp MsmpResponseSuccess) GetJSONRPC() string {
	return resp.JSONRPC
}

// IsSuccess 判断响应是否为成功响应
func (resp MsmpResponseFailure) IsSuccess() bool {
	return false
}

// GetID 获取请求ID
func (resp MsmpResponseFailure) GetID() int {
	return resp.ID
}

// GetJSONRPC 获取JSON-RPC版本
func (resp MsmpResponseFailure) GetJSONRPC() string {
	return resp.JSONRPC
}

// ParseResponse 解析响应数据，自动识别是成功还是失败响应
func ParseResponse(data []byte) (MsmpResponse, error) {
	// 创建一个临时结构来判断是否存在error字段
	temp := struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      int             `json:"id"`
		Result  json.RawMessage `json:"result,omitempty"`
		Error   json.RawMessage `json:"error,omitempty"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	// 根据是否存在error字段判断是成功还是失败响应
	if temp.Error != nil {
		// 失败响应
		failure := MsmpResponseFailure{
			JSONRPC: temp.JSONRPC,
			ID:      temp.ID,
		}

		if err := json.Unmarshal(temp.Error, &failure.Error); err != nil {
			return nil, err
		}

		return failure, nil
	} else {
		// 成功响应
		success := MsmpResponseSuccess{
			JSONRPC: temp.JSONRPC,
			ID:      temp.ID,
		}

		success.Result = temp.Result
		return success, nil
	}
}
