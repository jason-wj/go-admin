/*
 * @Author: lwnmengjing
 * @Date: 2021/6/10 3:39 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/6/10 3:39 下午
 */

package storage

import (
	config2 "go-admin/common/core/config"
	"log"

	"go-admin/common/core/sdk"
	"go-admin/common/core/sdk/pkg/captcha"
)

// Setup 配置storage组件
func Setup() {
	//4. 设置缓存
	cacheAdapter, err := config2.CacheConfig.Setup()
	if err != nil {
		log.Fatalf("cache setup error, %s\n", err.Error())
	}
	sdk.Runtime.SetCacheAdapter(cacheAdapter)
	//5. 设置验证码store
	captcha.SetStore(captcha.NewCacheStore(cacheAdapter, 600))

	//6. 设置队列
	if !config2.QueueConfig.Empty() {
		if q := sdk.Runtime.GetQueueAdapter(); q != nil {
			q.Shutdown()
		}
		queueAdapter, err := config2.QueueConfig.Setup()
		if err != nil {
			log.Fatalf("queue setup error, %s\n", err.Error())
		}
		sdk.Runtime.SetQueueAdapter(queueAdapter)
		defer func() {
			go queueAdapter.Run()
		}()
	}

	//7. 设置分布式锁
	if !config2.LockerConfig.Empty() {
		lockerAdapter, err := config2.LockerConfig.Setup()
		if err != nil {
			log.Fatalf("locker setup error, %s\n", err.Error())
		}
		sdk.Runtime.SetLockerAdapter(lockerAdapter)
	}
}
