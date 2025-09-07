package dto

type MsmpRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func NewMsmpRequest(id int, method string, param interface{}) MsmpRequest {
	return MsmpRequest{JSONRPC: "2.0", ID: id, Method: method, Params: param}
}
