// Copyright 2014 mqantserver Author. All Rights Reserved.
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
	"github.com/dyrkin/fsm"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"joysrv/comp"
	"math/rand"
	"reflect"
	"time"
)

var (
	SkeletonType []string = []string{"I", "O", "T", "J", "S", "L", "Z"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min) + min
}

type Table struct {
	room.QTable
	module          module.RPCModule
	players         []*Player
	state           string
	current_id      int
	current_frame   int //当前帧
	sync_frame      int //上一次同步数据的帧
	row             int
	col             int
	Speed           int64
	fsm             *fsm.FSM
	rule            Rule
	block_config    map[string][][]int
	tetris_template [][][]int
	grid            *Grid
	gameTypeInfo    *pb.S2S_Tetris_Create
	skeletonQueue   *SkeletonQueue
	heartbeat       *comp.Interval
}

func NewTable(module module.RPCModule, opts ...room.Option) (*Table, error) {
	this := &Table{
		module:        module,
		current_id:    0,
		current_frame: 0,
		sync_frame:    0,
		row:           15,
		col:           20,
		Speed:         1000,
	}

	opts = append(opts, room.TimeOut(60*5))                    //	如客户端操过5分钟未给房间发消息，房间将被回收
	opts = append(opts, room.Update(this.Update))              // 每30毫秒会调用一次该函数
	opts = append(opts, room.RunInterval(30*time.Millisecond)) // 时间周期设置为30毫秒
	//opts = append(opts, room.Capaciity(2500))
	opts = append(opts, room.NoFound(func(msg *room.QueueMsg) (value reflect.Value, e error) {
		return reflect.Zero(reflect.ValueOf("").Type()), errors.New("no found handler")
	}))
	opts = append(opts, room.SetRecoverHandle(func(msg *room.QueueMsg, err error) {
		log.Error("Recover %v Error: %v", msg.Func, err.Error())
	}))
	opts = append(opts, room.SetErrorHandle(func(msg *room.QueueMsg, err error) {
		log.Error("Error %v Error: %v", msg.Func, err.Error())
	}))
	this.OnInit(this, opts...)
	this.players = []*Player{
		NewPlayer(1, this.Speed),
		NewPlayer(2, this.Speed),
	}

	this.skeletonQueue = NewSkeletonQueue()
	this.Register("3", this.doJoin)
	this.Register("4", this.doExitRoom)
	this.Register("5", this.doSyncInfo)
	this.Register("6", this.doHD)
	this.Register("8", this.OperationSkeleton)
	this.heartbeat = comp.NewInterval(0, 1000) //1秒一次心跳
	return this, nil
}
