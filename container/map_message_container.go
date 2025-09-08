package container

import (
	"errors"
	"github.com/CycleZero/mc-msmp-go/dto"
)

type MapMessageContainer struct {
	WaitingMap map[int]*dto.MessagePair
	ReadyMap   map[int]*dto.MessagePair
}

func NewMapMessageContainer() *MapMessageContainer {
	return &MapMessageContainer{
		WaitingMap: make(map[int]*dto.MessagePair),
		ReadyMap:   make(map[int]*dto.MessagePair),
	}
}

func (m *MapMessageContainer) AddRequest(request *dto.MsmpRequest) error {
	id := request.ID
	_, exists := m.WaitingMap[id]
	if exists {
		return errors.New("duplicate request")
	}
	m.WaitingMap[id] = &dto.MessagePair{
		Id:       id,
		Request:  request,
		Response: nil,
	}
	return nil
}

func (m *MapMessageContainer) AddResponse(response dto.MsmpResponse) error {
	v, exists := m.WaitingMap[response.GetID()]
	if !exists {
		return errors.New("no waiting request")
	}
	v.Response = response
	m.ReadyMap[response.GetID()] = v
	delete(m.WaitingMap, response.GetID())
	return nil
}

func (m *MapMessageContainer) GetResponse(id int) (dto.MsmpResponse, error) {

	v, e := m.ReadyMap[id]
	if !e {
		return nil, errors.New("no response")
	}
	delete(m.ReadyMap, id)
	return v.Response, nil
}

func (m *MapMessageContainer) GetRequest(id int) (*dto.MsmpRequest, error) {
	v, e := m.WaitingMap[id]
	if !e {
		v, er := m.ReadyMap[id]
		if !er {
			return nil, errors.New("no request")
		}
		return v.Request, nil
	}
	return v.Request, nil
}

func (m *MapMessageContainer) GetResult(id int) (*dto.MessagePair, error) {
	v, e := m.ReadyMap[id]
	if !e {
		return nil, errors.New("no response")
	}
	delete(m.ReadyMap, id)
	return v, nil
}

func (m *MapMessageContainer) GetWaitingRequests() ([]*dto.MessagePair, error) {
	list := []*dto.MessagePair{}
	for _, v := range m.WaitingMap {
		list = append(list, v)
	}
	return list, nil
}

func (m *MapMessageContainer) GetWaitingNum() (int, error) {
	//var num int=0
	//for range m.WaitingMap {
	//	num++
	//}
	num := len(m.WaitingMap)

	return num, nil
}

func (m *MapMessageContainer) CancelRequest(id int) error {
	_, exists := m.WaitingMap[id]
	if !exists {
		return errors.New("no waiting request")
	}
	delete(m.WaitingMap, id)
	return nil
}

func (m *MapMessageContainer) AddRequestWithCallback(request *dto.MsmpRequest, callback func(*dto.MsmpRequest, dto.MsmpResponse)) error {
	id := request.ID
	_, exists := m.WaitingMap[id]
	if exists {
		return errors.New("duplicate request")
	}
	m.WaitingMap[id] = &dto.MessagePair{
		Id:       id,
		Request:  request,
		Response: nil,
		Callback: callback,
	}
	return nil
}
