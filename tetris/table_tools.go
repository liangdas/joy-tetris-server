package tetris

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
)

// GameTypeInfo 游戏房间初始化函数
func (table *Table) GameTypeInfo(req *pb.S2S_Tetris_Create) error {
	table.gameTypeInfo = req
	//初始化有限状态机
	table.fsm = table.initFSM()
	//初始化规则处理器
	table.rule = new(ClassicRule)
	//加载方块配置
	block_config, err := ReadBlocksFile(table.module.GetModuleSettings().Settings["BlocksPath"].(string))
	if err != nil {
		log.Error("ReadBlocksFile Fail", err.Error())
	} else {
		table.block_config = block_config
	}
	//加载快速模式下的地图模板数据
	tetris_template, err := ReadTemplateFile(table.module.GetModuleSettings().Settings["TemplatePath"].(string))
	if err != nil {
		log.Error("ReadBlocksFile Fail", err.Error())
	} else {
		table.tetris_template = tetris_template
	}
	if req.GameType == pb.TetrisGameType_CLASSIC {
		table.grid = NewGrid(table.row, table.col, table.col)
	} else if req.GameType == pb.TetrisGameType_QUICK {
		template := tetris_template[RandInt(0, len(tetris_template))]
		table.grid = NewGrid(table.row, table.col+len(template)*5, table.col)
		for range []int{0, 1, 2, 3} {
			template = table.grid.MergeTemplate(template, tetris_template[RandInt(0, len(tetris_template))])
		}
		table.grid.SetTemplate(table, template)
	} else {
		table.grid = NewGrid(table.row, table.col, table.col)
	}
	return nil
}

// SendToPlayer 给指定玩家发送消息
func (table *Table) SendToPlayer(session gate.Session, topic string, data interface{}) string {
	switch body := data.(type) {
	case []byte:
		return session.Send(topic, body)
	case string:
		return session.Send(topic, []byte(body))
	default:
		bb, err := json.Marshal(body)
		if err != nil {
			return err.Error()
		}
		return session.Send(topic, bb)
	}
	return ""
}

func (table *Table) Router() string {
	return table.Options().Router(table.TableId())
}

// GetBindPlayer 通过session获取已加入玩家位置信息
func (table *Table) GetBindPlayer(session gate.Session) room.BasePlayer {
	for _, player := range table.GetSeats() {
		if (player != nil) && (player.Session() != nil) {
			if player.Session().IsGuest() {
				if player.Session().GetSessionID() == session.GetSessionID() {
					player.OnRequest(session)
					return player
				}
			} else {
				if player.Session().GetUserID() == session.GetUserID() {
					player.OnRequest(session)
					return player
				}
			}

		}
	}
	return nil
}

// GetSpace 获取一个空位
func (table *Table) GetSpace(session gate.Session) *Player {
	for _, player := range table.players {
		if player.GetUserID() <= 0 {
			return player
		}
	}
	return nil
}

// FindByUID 通过uid获取位置信息
func (table *Table) FindByUID(user_id int64) *Player {
	for _, player := range table.players {
		if player.GetUserID() == user_id {
			return player
		}
	}
	return nil
}

// GetPlayers 获取所有玩家信息
func (this *Table) GetPlayers() []*Player {
	return this.players
}

// GetPlayerBySeatIndex 通过位置坐标获取位置信息
func (this *Table) GetPlayerBySeatIndex(SeatIndex int) *Player {
	for _, seat := range this.players {
		if seat.GetSeatIndex() == SeatIndex {
			return seat
		}
	}
	return nil
}

// PlayerNum 已加入游戏人数
func (table *Table) PlayerNum() int {
	num := 0
	for _, p := range table.players {
		if p.GetUserID() > 0 {
			num++
		}
	}
	return num
}

// GameStatus 游戏当前状态
func (table *Table) GameStatus() pb.TetrisGameStatus {
	switch table.state {
	case BeginState:
		return pb.TetrisGameStatus_BEGIN
	case ThisRoundOfSettlementState:
		return pb.TetrisGameStatus_GAMEOVER
	default:
		return pb.TetrisGameStatus_AWAIT
	}
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}
