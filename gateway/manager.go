package gateway

import (
	"strconv"
	"sync"
)

var manager sync.Map

// store 存储
func store(deviceId int64, c *Client) {
	manager.Store(deviceId, c)
	var key string = "userdevice:" + strconv.FormatInt(c.UserId, 10)
	redisClient.SAdd(key, c.DeviceId)
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
	value, ok := manager.Load(deviceId)
	if ok {
		UserId := value.(*Client).UserId
		var key string = "userdevice:" + strconv.FormatInt(UserId, 10)
		redisClient.SRem(key, deviceId)
	}
	manager.Delete(deviceId)
}
