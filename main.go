package main

import (
	"fmt"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/log/beego"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/rpc/util"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"joysrv/account"
	"joysrv/gate"
	"joysrv/tetris"
	"joysrv/webapp"
	"net/http"
	"time"
)

const timeFormat = "20060102T150405Z"

func main() {
	go func() {
		//http://127.0.0.1:7070/debug/pprof/
		http.ListenAndServe("0.0.0.0:7070", nil)
	}()
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	//nc, err := nats.Connect("nats://192.168.20.119:4222")
	if err != nil {
		log.Error("nats error %v", err)
		return
	}
	app := mqant.CreateApp(
		module.Debug(true), //是否开启debug模式
		module.Nats(nc),    //指定nats rpc
		//module.WorkDir("/work/go/joy-tetris-server"),
		//module.Configure("/work/go/joy-tetris-server/bin/conf/server.json"), // 配置
		module.ProcessID("development"), //模块组ID
	)
	_ = app.SetProtocolMarshal(func(Trace string, Result interface{}, Error string) (module.ProtocolMarshal, string) {
		var result []byte
		if Result != nil {
			//内容不为空,尝试转为[]byte
			switch v2 := Result.(type) {
			case module.ProtocolMarshal:
				result = v2.GetData()
			default:
				_, r, err := argsutil.ArgsTypeAnd2Bytes(app, Result)
				if err != nil {
					Error = err.Error()
				}
				result = r
			}
		}
		r := &pb.S2C_Response{
			Error:  *proto.String(Error),
			Trace:  *proto.String(Trace),
			Result: result,
		}
		b, err := proto.Marshal(r)
		if err == nil {
			//解析得到[]byte后用NewProtocolMarshal封装为module.ProtocolMarshal
			return app.NewProtocolMarshal(b), ""
		} else {
			return nil, err.Error()
		}
	})
_:
	app.OnConfigurationLoaded(func(app module.App) {
		fmt.Println(time.Now().UTC().Format(timeFormat))
	})
_:
	app.OnStartup(func(app module.App) {
		log.LogBeego().SetFormatFunc(logs.DefineErrorLogFunc(app.GetProcessID(), 4))
	})
_:
	app.Run(
		account.Module(),
		mgate.Module(),
		tetris.Module(),
		webapp.Module(),
	)
}
