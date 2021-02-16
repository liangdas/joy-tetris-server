// Copyright 2014 loolgame Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tetris

import (
	"errors"
	"fmt"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"google.golang.org/protobuf/proto"
)

var Module = func() module.Module {
	this := new(tetris)
	return this
}

type tetris struct {
	basemodule.BaseModule
	room    *room.Room
	proTime int64
	gameId  int
}

func (self *tetris) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "tetris"
}
func (self *tetris) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *tetris) GetFullServerId() string {
	return self.GetType() + "@" + self.GetServerId()
}

func (self *tetris) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.gameId = 13
	self.room = room.NewRoom(self.GetApp())
	self.GetServer().Register(pb.S2STetrisCreate, self.createRoomDebug)
	self.GetServer().Register(pb.C2STetris, self.action)
}

func (self *tetris) Run(closeSig chan bool) {

}

func (self *tetris) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

func (self *tetris) NewTetrisTable(module module.App, tableId string) (room.BaseTable, error) {
	table, err := NewTable(
		self,
		room.TableId(tableId),
		room.Router(func(TableId string) string {
			return fmt.Sprintf("%v://%v/room", self.GetType(), self.GetServer().ID())
		}),
		room.DestroyCallbacks(func(table room.BaseTable) error {
			log.Info("回收了房间: %v", table.TableId())
			_ = self.room.DestroyTable(table.TableId())
			return nil
		}),
	)
	return table, err
}

// createRoomDebug 创建游戏房间
func (self *tetris) createRoomDebug(session gate.Session, req *pb.S2S_Tetris_Create) (*pb.S2R_Tetris_Create, error) {
	tableid := req.GetRoomId()
	table_id := fmt.Sprintf("%v", tableid)
	table := self.room.GetTable(table_id)
	if table == nil {
		table, err := self.room.CreateById(self.GetApp(), table_id, self.NewTetrisTable)
		if err != nil {
			return nil, err
		}
		table.(*Table).GameTypeInfo(req)
		table.Run()
	} else {
		return nil, errors.New("room_already_exists")
	}
	errstr := session.SetPush(self.GetType(), table_id)
	if errstr != "" {
		log.Error("SetPush ludo server_id err", errstr)
	}
	return &pb.S2R_Tetris_Create{
		GameType:  req.GameType,
		PlayerNum: *proto.Int64(req.GetPlayerNum()),
		Gear:      *proto.Int64(req.GetGear()),
		UseProps:  *proto.Bool(req.GetUseProps()),
		UserList:  req.GetUserList(),
		RoomId:    *proto.String(table_id),
		ShortId:   *proto.String(table_id),
		Private:   *proto.Bool(req.GetPrivate()),
		Version:   *proto.String(req.GetVersion()),
		Owner:     *proto.Int64(req.GetOwner()),
	}, nil
}

// action 游戏指令接收函数
func (self *tetris) action(session gate.Session, req *pb.C2S_Tetris) (r string, err error) {
	table_id := req.RoomId
	table := self.room.GetTable(table_id)
	if table == nil {
		return "", errors.New("nofound room")
	}
	erro := table.PutQueue(req.MsgType, session, req)
	if erro != nil {
		return "", erro
	}
	return "success", nil
}
