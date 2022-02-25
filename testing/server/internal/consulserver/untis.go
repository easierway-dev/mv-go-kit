package consulserver

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
	"unsafe"
)

func getTomlConfig(ops *Ops, value interface{}) error {
	// 获取toml配置文件中的值
	pair, err := GetValue(ops)
	if err != nil {
		return err
	}
	if pair == nil {
		return KvNotFound
	}
	// 将配置文件的值与consulConfig进行绑定
	if _, err = toml.Decode(*(*string)(unsafe.Pointer(&pair.Value)), value); err != nil {
		return err
	}
	// MD5加密
	value.(*ServersConfig).hashTag = MD5(string(pair.Value))
	return nil
}

func GetValue(ops *Ops) (*api.KVPair, error) {
	config := api.DefaultConfig()
	config.Address = ops.Address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	kv := client.KV()
	if kv == nil {
		return nil, GetConsulKvFailed
	}
	pair, _, err := kv.Get(ops.Path, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, KvNotFound
	}
	return pair, nil
}
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func timedTask(sm *ServerManager) {
	var sy sync.WaitGroup
	ticker := time.NewTicker(time.Second)
	sy.Add(1)
	for {
		select {
		case <-ticker.C:
			defer sy.Done()

		}
	}
	sy.Wait()
}

func RunTask(sm *ServerManager) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			TimeTask(sm)
		}
	}
}

func TimeTask(sm *ServerManager) {
	fmt.Println("定时任务开始:")
	sc, _ := FromConsulConfig("127.0.0.1:8500", CONSULKEY)
	sm.sc = sc
}
