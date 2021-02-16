// Copyright 2014 sdkgame Author. All Rights Reserved.
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
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
)

// doHD 心跳
func (table *Table) doHD(session gate.Session, req *pb.C2S_Tetris) (err error) {
	player := table.FindPlayer(session)
	if player == nil {
		return errors.New("no join")
	}
	player.OnRequest(session)
	return nil
}

// doJoin 客户端加入房间指令
func (table *Table) doJoin(session gate.Session, req *pb.C2S_Tetris) (err error) {
	bindplayer := table.GetBindPlayer(session)
	if bindplayer != nil {
		bindplayer.OnRequest(session)
		if table.grid != nil {
			table.grid.ResetData()
		}
		table.S2C_GameStatusChangeBroadcast()
		table.S2C_PlayerStatusChangeBroadcast(bindplayer.(*Player), pb.S2C_PlayerStatusChangeBroadcast_ONLINE)
		return nil
	} else {
		player := table.FindByUID(session.GetUserIdInt64())
		if player != nil {
			player.OnRequest(session)
			if table.grid != nil {
				table.grid.ResetData()
			}
			table.S2C_GameStatusChangeBroadcast()
			table.S2C_PlayerStatusChangeBroadcast(player, pb.S2C_PlayerStatusChangeBroadcast_EXIT)
			return nil
		}
	}
	player := table.GetSpace(session)
	//player := table.FindByUID(session.GetUserIdInt64())
	if player == nil {
		_ = table.SendToPlayer(session, "/room/exit", []byte("加入游戏失败!"))
		return nil
	}
	player.UserID = session.GetUserIdInt64()
	player.Bind(session)
	player.OnRequest(session)
	if table.grid != nil {
		table.grid.ResetData()
	}
	table.S2C_GameStatusChangeBroadcast()
	table.S2C_PlayerStatusChangeBroadcast(player, pb.S2C_PlayerStatusChangeBroadcast_ENTER)
	return nil
}

// doSyncInfo 客户端申请同步全量数据
func (table *Table) doSyncInfo(session gate.Session, req *pb.C2S_Tetris) (err error) {
	playerImp := table.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*Player)
		player.OnRequest(session)
		table.S2C_GameStatusChangeBroadcast()
		if table.grid != nil {
			table.grid.ResetData()
		}
		return nil
	} else {
		log.Error("doSyncInfo unbind")
	}
	return err
}

// doSyncInfo 客户端退出房间
func (table *Table) doExitRoom(session gate.Session, req *pb.C2S_Tetris) (err error) {
	playerImp := table.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*Player)
		player.OnRequest(session)
		//解绑座位，后续不能再给玩家发送消息，否则会造成玩家其他游戏局内的数据干扰
		player.UnBind()
		table.S2C_GameStatusChangeBroadcast()
		table.S2C_PlayerStatusChangeBroadcast(player, pb.S2C_PlayerStatusChangeBroadcast_EXIT)
	} else {
		//如果未加入游戏,则直接通过uid找到玩家
		player := table.FindByUID(session.GetUserIdInt64())
		if player != nil {
			table.S2C_GameStatusChangeBroadcast()
			table.S2C_PlayerStatusChangeBroadcast(player, pb.S2C_PlayerStatusChangeBroadcast_EXIT)
		}
	}
	return err
}

// OperationSkeleton 客户端操作方块指令
func (this *Table) OperationSkeleton(session gate.Session, req *pb.C2S_Tetris) {
	playerImp := this.GetBindPlayer(session)
	if playerImp != nil {
		player := playerImp.(*Player)
		player.OnRequest(session)
		op := req.GetPlayerOperationSkeletonC2S()
		//TODO 通过限定方块Frame来防止误操作
		this.skeletonQueue.ExecutionSkeleton(this, player, player.GetSeatIndex(), int(op.GetFrame()), op.GetOpcode(), nil)
	}
}
