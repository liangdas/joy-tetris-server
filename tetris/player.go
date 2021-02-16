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
	"fmt"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant-modules/room"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	room.BasePlayerImp
	UserID       int64
	SeatIndex    int
	Robot        bool  //是否是机器人
	Score        int   //获得的分数
	Acceleration int   //方块快速下滑的速度等级 0 50ms  1  50ms 2次
	Speed        int64 //正常下降的速度
}

func NewPlayer(SeatIndex int, Speed int64) *Player {
	this := new(Player)
	this.SeatIndex = SeatIndex
	this.Speed = Speed
	this.Acceleration = 1
	this.Robot = false
	return this
}

func (this *Player) IsRobot() bool {
	return this.Robot
}

func (this *Player) GetSeatIndex() int {
	return this.SeatIndex
}
func (this *Player) GetUserID() int64 {
	return this.UserID
}
func (this *Player) SetSeatIndex(SeatIndex int) {
	this.SeatIndex = SeatIndex
}
func (this *Player) GetScore() int {
	return this.Score
}
func (this *Player) SetScore(Score int) {
	this.Score = Score
}
func (this *Player) GetAcceleration() int {
	return this.Acceleration
}
func (this *Player) SetAcceleration(Acceleration int) {
	this.Acceleration = Acceleration
}

func (this *Player) AccelerationSpeed() int64 {
	if this.Acceleration == 0 {
		return 50
	} else if this.Acceleration == 1 {
		return 25
	} else if this.Acceleration == 2 {
		return 25
	}
	return 100
}
func (this *Player) GetSpeed() int64 {
	return this.Speed
}
func (this *Player) SetSpeed(Speed int64) {
	this.Speed = Speed
}

func (lp *Player) ProtoBuffer() *pb.TetrisPlayerInfo {
	return &pb.TetrisPlayerInfo{
		SeatIndex: *proto.Uint32(uint32(lp.SeatIndex)),
		UserId:    *proto.Int64(lp.GetUserID()),
		Score:     *proto.Uint32(uint32(lp.Score)),
		InGame:    *proto.Bool(true),
		Online:    *proto.Bool(lp.IsBind()),
		Auto:      *proto.Bool(false),
		NickName:  *proto.String(fmt.Sprintf("%v", lp.GetUserID())),
		IconUrl:   *proto.String(""),
		IconFrame: *proto.String(""),
	}
}
