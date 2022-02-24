package consulserver

import (
	"context"
	"fmt"
	"time"
)

type ServerManager struct {
	sc      *ServersConfig
	servers map[int]*Server // key: port
	hashTag string          // current config tag
	status  bool
}

func NewServerManager() *ServerManager {
	fmt.Println("Create ServerManager")
	// sc := GetServersConfigFromConsul()
	sc := GetServersConfigFromLocal()
	//sc, _ := FromConsulConfig("47.252.4.203:8500", "/jianjilong")
	sm, err := NewServerManagerWithConfig(sc)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return sm
}

func (sm *ServerManager) Serve(ctx context.Context) {
	// todo: clear all service node
	select {
	case <-ctx.Done():
		return
	}
}
func NewServerManagerWithConfig(sc *ServersConfig) (*ServerManager, error) {
	fmt.Println("Create ServerManager From Config")
	sm := &ServerManager{
		sc:      sc,
		servers: make(map[int]*Server),
	}
	err := sm.manage()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return sm, nil
}
func (sm *ServerManager) manage() error {
	sm.Clear()
	ticker := time.NewTicker(time.Second)
	sm.sync()
	go func() {
		for {
			<-ticker.C
			sm.sync()
		}
	}()
	return nil
}
func (sm *ServerManager) Clear() {
	fmt.Println("deregister all services", SERVICE)
}

func (sm *ServerManager) sync() {
	/*
		将consul的配置同步到真正的服务上
	*/
	fmt.Println("start sync")
	if sm.hashTag == sm.sc.hashTag {
		// 配置没变，啥也不干
		return
	}
	// 配置有问题, 啥也不干
	serverConfigs := sm.sc.GetServerConfigs()
	if len(serverConfigs) == 0 {
		sm.status = false
		return
	}

	// create or update and register consulserver
	for port, serverProperty := range serverConfigs {
		// 不在servers这个map中，创建一个server
		if _, ok := sm.servers[port]; !ok {
			sm.servers[port] = NewServer(port)
		}
		server := sm.servers[port]
		if err := server.applyProperty(serverProperty); err != nil {
			fmt.Println(err)
			sm.status = false
			return
		}
	}
	// remove and deregister consulserver
	for port, server := range sm.servers {
		if _, ok := serverConfigs[port]; !ok {
			server.destroy()
			delete(sm.servers, port)
		}
	}
	sm.status = true
	sm.hashTag = sm.sc.hashTag
}
