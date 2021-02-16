/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package account

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

var Module = func() module.Module {
	user := new(Account)
	user.endCond = false
	return user
}

type JwtUserClaims struct {
	UserId string `json:"user_id"`
	IP     string `json:"ip"`
	TS     int64  `json:"ts"`
	Sign   string `json:"sign"`
	jwt.StandardClaims
}
type Account struct {
	basemodule.BaseModule
	endCond       bool
	recommend_ch  chan *nats.Msg
	recommend_sub *nats.Subscription
	usertable     *room.Room
}

func (self *Account) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "account"
}
func (self *Account) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Account) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.usertable = room.NewRoom(app)
	self.GetServer().RegisterGO(pb.C2SLoginDEBUG, self.LoginWithUserId)

}

func (self *Account) Run(closeSig chan bool) {
	<-closeSig
}

func (self *Account) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

// LoginWithUserId 长连接验证绑定UID
func (acc *Account) LoginWithUserId(session gate.Session, req *pb.C2S_Login_DEBUG) (*pb.S2C_Login, error) {
	//TODO 验证身份信息合法性
	if session != nil {
		errstr := session.Bind(fmt.Sprintf("%v", req.UserId))
		if errstr != "" {
			return nil, errors.New(errstr)
		}
	}
	return &pb.S2C_Login{
		UserId: *proto.Int64(req.UserId),
		Nick:   "",
		Avatar: *proto.String(""),
	}, nil
}
