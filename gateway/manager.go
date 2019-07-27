package gateway

import "sync"

var manager sync.Map

// store 存储
func store(deviceId int64, c *Client) {
	manager.Store(deviceId, c)
}

// load 获取
func load(deviceId int64) *Client {
	value, ok := manager.Load(deviceId)
	if ok {
		return value.(*Client)
	}
	return nil
}

// delete 删除
func delete(deviceId int64) {
	manager.Delete(deviceId)
}
