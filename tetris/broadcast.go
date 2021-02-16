package tetris

import (
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"google.golang.org/protobuf/proto"
)

// S2C_PlayerStatusChangeBroadcast 玩家进入房间广播
func (table *Table) S2C_PlayerStatusChangeBroadcast(player *Player, status pb.S2C_PlayerStatusChangeBroadcast_PlayerStatusType) error {
	playerInfo := []*pb.TetrisPlayerInfo{}
	for _, p := range table.players {
		if p.GetUserID() > 0 {
			playerInfo = append(playerInfo, p.ProtoBuffer())
		}
	}
	p := &pb.S2C_PlayerStatusChangeBroadcast{
		UserId:     *proto.Int64(player.GetUserID()),
		SeatIndex:  *proto.Uint32(uint32(player.SeatIndex)),
		Status:     status,
		PlayerInfo: playerInfo,
	}

	return table.Broadcast(p)
}

// S2C_GameStatusChangeBroadcast 房间状态变更广播
func (table *Table) S2C_GameStatusChangeBroadcast() error {
	playerInfo := []*pb.TetrisPlayerInfo{}
	for _, p := range table.players {
		if p.GetUserID() > 0 {
			playerInfo = append(playerInfo, p.ProtoBuffer())
		}
	}
	p := &pb.S2C_GameStatusChangeBroadcast{
		RoomInfo: &pb.TetrisRoomInfo{
			Status:    table.GameStatus(),
			PlayerNum: *proto.Int64(int64(table.gameTypeInfo.PlayerNum)),
			RoomType:  table.gameTypeInfo.GameType,
		},
		PlayerInfo: playerInfo,
	}

	return table.Broadcast(p)
}

// Broadcast 广播消息
func (table *Table) Broadcast(p interface{}) error {
	var resp *pb.S2C_Tetris
	switch p.(type) {
	case *pb.S2C_PlayerStatusChangeBroadcast:
		resp = &pb.S2C_Tetris{
			MsgId:                          *proto.Uint32(uint32(table.current_frame)),
			MsgType:                        *proto.Uint32(4),
			RoomId:                         *proto.String(table.TableId()),
			PlayerStatusChangeS2CBroadcast: p.(*pb.S2C_PlayerStatusChangeBroadcast),
		}
	case *pb.S2C_GameStatusChangeBroadcast:
		resp = &pb.S2C_Tetris{
			MsgId:                        *proto.Uint32(uint32(table.current_frame)),
			MsgType:                      *proto.Uint32(5),
			RoomId:                       *proto.String(table.TableId()),
			GameStatusChangeS2CBroadcast: p.(*pb.S2C_GameStatusChangeBroadcast),
		}
	}
	b, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	return table.NotifyRealMsg("/s2c_tetris/", b)
}
