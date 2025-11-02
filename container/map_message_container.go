package container

import (
	"errors"
	"github.com/CycleZero/mc-msmp-go/dto"
	"sync"
)

type MapMessageContainer struct {
	WaitingMap map[int]*dto.MessagePair
	ReadyMap   map[int]*dto.MessagePair
	mutex      sync.RWMutex
}

func NewMapMessageContainer() *MapMessageContainer {
	return &MapMessageContainer{
		WaitingMap: make(map[int]*dto.MessagePair),
		ReadyMap:   make(map[int]*dto.MessagePair),
	}
}
func (m *MapMessageContainer) NewResponse(r dto.MsmpResponse) (*dto.MessagePair, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, e := m.WaitingMap[r.GetID()]
	if !e {
		return nil, errors.New("no waiting request")
	}
	v.Response = r
	delete(m.WaitingMap, r.GetID())
	return v, nil

}

func (m *MapMessageContainer) AddRequest(request *dto.MsmpRequest) error {
	id := request.ID
	m.mutex.RLock()
	_, exists := m.WaitingMap[id]
	m.mutex.RUnlock()
	if exists {
		return errors.New("duplicate request")
	}
	m.mutex.Lock()
	m.WaitingMap[id] = &dto.MessagePair{
		Id:       id,
		Request:  request,
		Response: nil,
	}
	m.mutex.Unlock()
	return nil
}

// AddResponse Deprecated
func (m *MapMessageContainer) AddResponse(response dto.MsmpResponse) error {
	m.mutex.Lock()
	v, exists := m.WaitingMap[response.GetID()]
	if !exists {
		return errors.New("no waiting request")
	}
	v.Response = response
	//m.ReadyMap[response.GetID()] = v
	if v.Callback != nil {
		go v.Callback(v.Request, v.Response)
	}
	delete(m.WaitingMap, response.GetID())
	m.mutex.Unlock()
	return nil
}

// Deprecated
func (m *MapMessageContainer) GetResponse(id int) (dto.MsmpResponse, error) {

	v, e := m.ReadyMap[id]
	if !e {
		return nil, errors.New("no response")
	}
	delete(m.ReadyMap, id)
	return v.Response, nil
}

func (m *MapMessageContainer) GetRequest(id int) (*dto.MsmpRequest, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
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

// Deprecated
func (m *MapMessageContainer) GetResult(id int) (*dto.MessagePair, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, e := m.ReadyMap[id]
	if !e {
		return nil, errors.New("no response")
	}
	delete(m.ReadyMap, id)
	return v, nil
}

func (m *MapMessageContainer) GetWaitingRequests() ([]*dto.MessagePair, error) {
	list := []*dto.MessagePair{}
	m.mutex.RLock()
	defer m.mutex.RUnlock()
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	num := len(m.WaitingMap)

	return num, nil
}

func (m *MapMessageContainer) CancelRequest(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, exists := m.WaitingMap[id]
	if !exists {
		return errors.New("no waiting request")
	}
	delete(m.WaitingMap, id)
	return nil
}

func (m *MapMessageContainer) AddRequestWithHandler(request *dto.MsmpRequest, callback func(*dto.MsmpRequest, dto.MsmpResponse)) error {
	id := request.ID
	m.mutex.Lock()
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
	m.mutex.Unlock()
	return nil
}
