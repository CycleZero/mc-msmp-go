package iface

import "github.com/CycleZero/mc-msmp-go/dto"

type MessageContainer interface {
	AddRequest(request *dto.MsmpRequest) error
	AddRequestWithHandler(request *dto.MsmpRequest, handler func(*dto.MsmpRequest, dto.MsmpResponse)) error
	//AddResponse(response dto.MsmpResponse) error
	//GetResponse(id int) (dto.MsmpResponse, error)
	GetRequest(id int) (*dto.MsmpRequest, error)
	//GetResult(id int) (*dto.MessagePair, error)
	NewResponse(r dto.MsmpResponse) (*dto.MessagePair, error)
	GetWaitingRequests() ([]*dto.MessagePair, error)
	GetWaitingNum() (int, error)
	CancelRequest(id int) error
}
