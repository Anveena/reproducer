package quic

import (
	"sync"
	"time"
)

type client struct {
	wg *sync.WaitGroup
}
type manager struct {
	clientDicLock *sync.RWMutex
	clientDic     map[uint32][]*client
}

var _manager *manager

func init() {
	_manager = &manager{
		clientDicLock: &sync.RWMutex{},
		clientDic:     make(map[uint32][]*client),
	}
}

func getManager() *manager {
	return _manager
}

func (m *manager) join(generation, total uint32, wg *sync.WaitGroup) {
	m.clientDicLock.Lock()
	defer m.clientDicLock.Unlock()
	if _, ok := m.clientDic[generation]; !ok {
		m.clientDic[generation] = []*client{}
		time.AfterFunc(time.Second*10, func() {
			m.clientDicLock.Lock()
			defer m.clientDicLock.Unlock()
			list := m.clientDic[generation]
			delete(m.clientDic, generation)
			if len(list) == 0 {
				return
			}
			for _, client := range list {
				client.wg.Done()
			}
		})
	}
	m.clientDic[generation] = append(m.clientDic[generation], &client{wg: wg})
	println("join", generation, total, len(m.clientDic[generation]))
	if len(m.clientDic[generation]) == int(total) {
		for _, client := range m.clientDic[generation] {
			client.wg.Done()
		}
		delete(m.clientDic, generation)
		println("clean up")
	}
}
