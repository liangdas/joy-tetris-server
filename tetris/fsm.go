package tetris

import (
	"github.com/dyrkin/fsm"
	"joysrv/component"
	"math/rand"
	"time"
)

//states
const InitialState = "InitialState"                        //初始化
const AwaitState = "Await"                                 //等待开始游戏
const BeginState = "Begin"                                 //开始游戏
const ThisRoundOfSettlementState = "ThisRoundOfSettlement" //本轮结算

//messages
type Transfer struct {
	source chan int
	target chan int
	amount int
}

const GoForward = "GoForward"
const DoJoin = "DoJoin"

type TransferData struct {
	event   string
	message interface{}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// initFSM 初始化游戏的逻辑管理有限状态机
func (table *Table) initFSM() *fsm.FSM {
	wt := fsm.NewFSM()

	wt.SetDefaultHandler(func(event *fsm.Event) *fsm.NextState {
		//啥也不处理
		return wt.Stay()
	})

	wt.StartWith(InitialState, table)

	wt.When(InitialState)(table.InitialState)

	wt.When(AwaitState)(table.AwaitState)

	wt.When(BeginState)(table.BeginState)

	wt.When(ThisRoundOfSettlementState)(
		table.GameSettlementState)

	return wt
}

func (table *Table) fsmEvent(event string, message interface{}) error {
	table.fsm.Send(&TransferData{event: event, message: message})
	return nil
}

func (table *Table) InitialState(event *fsm.Event) *fsm.NextState {
	if table.state != InitialState {
		table.state = InitialState
		table.S2C_GameStatusChangeBroadcast()
	}
	return table.fsm.Goto(AwaitState)
}

func (table *Table) AwaitState(event *fsm.Event) *fsm.NextState {
	if table.state != AwaitState {
		table.state = AwaitState
		table.S2C_GameStatusChangeBroadcast()
	}
	wt := table.fsm
	message, dataOk := event.Message.(*TransferData)
	if dataOk {
		switch message.event {
		case GoForward, DoJoin:
			if table.rule.StartGame(table) {

				return wt.Goto(BeginState)
			} else {
				return wt.Stay()
			}
		}
	}
	return wt.DefaultHandler()(event)
}

func (table *Table) BeginState(event *fsm.Event) *fsm.NextState {
	if table.state != BeginState {
		table.state = BeginState
		table.S2C_GameStatusChangeBroadcast()
	}
	wt := table.fsm
	message, dataOk := event.Message.(*TransferData)
	if dataOk {
		switch message.event {
		case GoForward:
			if table.rule.EndOfGame(table) {
				//游戏已结束，就去结算逻辑判断一下
				return wt.Goto(ThisRoundOfSettlementState)
			}
			table.current_frame++
			for _, player := range table.GetPlayers() {
				if player.GetUserID() != 0 && table.skeletonQueue.GetSkeleton(player.GetSeatIndex()) == nil {
					Type := SkeletonType[RandInt(0, len(SkeletonType))]
					skeleton := NewSkeleton(Type, player.GetSeatIndex(), RandInt(0, 4), player.GetSeatIndex()*4+1, 0, table.current_frame)
					table.skeletonQueue.AddSkeleton(skeleton) //Type string,SeatIndex,Index ,startx,starty
					table.grid.UpdateData()
					skeleton.Start(table)
				}
			}
			table.skeletonQueue.Update(table)
			table.grid.Update(table)
			data, ty, err := table.grid.Marshal(table)
			if err == nil {
				if ty == component.FULL {
					g, err := GzipEncode(data)
					if err == nil {
						//log.Info("stores,GZIP %v  %v", len(g), len(data))
						_ = table.NotifyRealMsg("/tetris/grid/data/", g)
					}

				} else if ty == component.PATCH {
					//log.Info("table,PATCH %v", len(data))
					_ = table.NotifyRealMsg("/tetris/grid/patch/", data)
				}
			}
		}
	}
	return wt.DefaultHandler()(event)
}

func (table *Table) GameSettlementState(event *fsm.Event) *fsm.NextState {
	if table.state != ThisRoundOfSettlementState {
		table.state = ThisRoundOfSettlementState
		table.S2C_GameStatusChangeBroadcast()
	}
	wt := table.fsm
	message, dataOk := event.Message.(*TransferData)
	if dataOk {
		switch message.event {
		case GoForward:
			table.Finish()
			return wt.Stay()
		}
	}
	return wt.DefaultHandler()(event)
}
