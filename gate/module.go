/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package mgate

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/gate/base"
	"github.com/liangdas/mqant/gate/uriroute"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"net/url"
	"time"
)

type LocalUserData struct {
	ProHeartbeatTime           time.Time
	ProSyncOnlineHeartbeatTime time.Time
}

var Module = func() module.Module {
	this := new(Gate)
	return this
}

type Gate struct {
	basegate.Gate //继承
	RedisUrl      string
	Route         *uriroute.URIRoute
}

func (this *Gate) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "gate"
}
func (this *Gate) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (this *Gate) OnInit(app module.App, settings *conf.ModuleSettings) {

	route := uriroute.NewURIRoute(this,
		uriroute.Selector(this.Selector),
		uriroute.DataParsing(func(topic string, u *url.URL, msg []byte) (bean interface{}, err error) {
			return
		}),
		uriroute.CallTimeOut(3*time.Second),
	)
	//注意这里一定要用 gate.Gate 而不是 module.BaseModule
	this.Gate.OnInit(this, app, settings,
		gate.Heartbeat(time.Second*5),
		gate.BufSize(2048*20),
		gate.SetRouteHandler(route),
		gate.SetSessionLearner(this),
		gate.SetStorageHandler(this),
	)
	this.Gate.SetJudgeGuest(func(session gate.Session) bool {
		if session.GetUserID() == "" {
			return true
		}
		return false
	})
}

//当连接建立  并且MQTT协议握手成功
func (this *Gate) Connect(session gate.Session) {
	log.Info("客户端建立了链接")
	_ = session.SetLocalUserData(&LocalUserData{
		ProHeartbeatTime:           time.Now(),
		ProSyncOnlineHeartbeatTime: time.Now(),
	})
}

//当连接关闭	或者客户端主动发送MQTT DisConnect命令 ,这个函数中Session无法再继续后续的设置操作，只能读取部分配置内容了
func (this *Gate) DisConnect(session gate.Session) {
	log.Info("客户端断开了链接")
}

func (this *Gate) OnRoute(session gate.Session, topic string, msg []byte) (bool, interface{}, error) {
	return this.Route.OnRoute(session, topic, msg)
}

/**
存储用户的Session信息
Session Bind Userid以后每次设置 settings都会调用一次Storage
*/
func (this *Gate) Storage(session gate.Session) (err error) {
	return nil
}

/**
强制删除Session信息
*/
func (this *Gate) Delete(session gate.Session) (err error) {
	return
}

/**
获取用户Session信息
用户登录以后会调用Query获取最新信息
*/
func (this *Gate) Query(Userid string) ([]byte, error) {
	return nil, nil
}

/**
用户心跳,一般用户在线时60s发送一次
可以用来延长Session信息过期时间
*/
func (this *Gate) Heartbeat(session gate.Session) {
	userdata, ok := session.LocalUserData().(*LocalUserData)
	if ok {
		onlinehbtime := time.Now().Sub(userdata.ProSyncOnlineHeartbeatTime)
		if onlinehbtime.Seconds() > 60*2 {
			//每2s同步一次
		}
	}
}
