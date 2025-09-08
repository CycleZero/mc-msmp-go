package iface

import "github.com/CycleZero/mc-msmp-go/dto"

type MessageContainer interface {
	AddRequest(request *dto.MsmpRequest) error
	AddRequestWithCallback(request *dto.MsmpRequest, callback func(*dto.MsmpRequest, dto.MsmpResponse)) error
	AddResponse(response dto.MsmpResponse) error
	GetResponse(id int) (dto.MsmpResponse, error)
	GetRequest(id int) (*dto.MsmpRequest, error)
	GetResult(id int) (*dto.MessagePair, error)
	GetWaitingRequests() ([]*dto.MessagePair, error)
	GetWaitingNum() (int, error)
	CancelRequest(id int) error
}
