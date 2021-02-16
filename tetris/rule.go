package tetris

type Rule interface {
	//判断是否具备开始游戏条件
	StartGame(table *Table) bool
	//判断是否获得了胜利
	Win(table *Table, player *Player) bool
	//游戏结束判断
	EndOfGame(table *Table) bool

	//结算
	Settlement(table *Table) error
}

//经典模式的规则
type ClassicRule struct {
}

func (rule *ClassicRule) StartGame(table *Table) bool {
	ready := true
	if table.PlayerNum() < int(table.gameTypeInfo.PlayerNum) {
		ready = false
	}
	if !ready {
		return false
	}
	return true
}

func (rule *ClassicRule) Win(table *Table, player *Player) bool {
	end := true
	return end
}

func (rule *ClassicRule) EndOfGame(table *Table) bool {

	return table.grid.GameOver()
}

func (rule *ClassicRule) Settlement(table *Table) error {

	return nil
}
