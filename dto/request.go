package dto

type MsmpRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

type MessagePair struct {
	Id       int
	Request  *MsmpRequest
	Response MsmpResponse
	Callback func(request *MsmpRequest, response MsmpResponse)
}

func NewMsmpRequest(id int, method string, param interface{}) MsmpRequest {
	return MsmpRequest{JSONRPC: "2.0", ID: id, Method: method, Params: []interface{}{param}}
}
