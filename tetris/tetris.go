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
	"github.com/golang/protobuf/proto"
	pb "github.com/liangdas/joy-tetris-protobuf/golang/message"
	"github.com/liangdas/mqant-modules/component"
	"github.com/pkg/errors"
	"joysrv/comp"
	"math"
)

var (
	NORTH = [2]int{0, -1} //北
	SOUTH = [2]int{0, 1}  //南
	WEST  = [2]int{-1, 0} //西
	EAST  = [2]int{1, 0}  //东
)

type Point struct {
	x int
	y int
}

func (this *Point) X() int {
	return this.x
}

func (this *Point) Y() int {
	return this.y
}

func NewSkeleton(Type string, SeatIndex, Index, startx, starty, frame int) *Skeleton {
	this := &Skeleton{
		ttype:     Type,
		index:     Index,
		frame:     frame,
		seatIndex: SeatIndex,
		started:   false,
		point: &Point{
			x: startx,
			y: starty,
		},
		interval: comp.NewInterval(0, 1000),
	}
	return this
}

// 方块结构体，4*4矩阵
//0,0,0,0,
//1,1,1,1,
//0,0,0,0,
//0,0,0,0
type Skeleton struct {
	ttype     string  //方块类型，跟配置中对应 //类型 I O T J S L Z
	index     int     //变换索引 0 - 3
	frame     int     //创建这个块的帧数，可以用来防抖
	seatIndex int     // 方块所属玩家
	point     *Point  //方块当前所在游戏画布中的坐标
	lattice   [][]int //点阵4*4
	started   bool
	interval  *comp.Interval
}

func (this *Skeleton) GetSeatIndex() int {
	return this.seatIndex
}

func (this *Skeleton) Start(table *Table) {
	if this.started {
		return
	}

	var err error = nil
	for err == nil {
		err = table.grid.Move(table, this, NORTH) //移动到最顶上
	}
	this.started = true
}

func (this *Skeleton) Started() bool {
	return this.started
}

func (this *Skeleton) Type() string {
	return this.ttype
}

func (this *Skeleton) SetType(ttype string) {
	this.ttype = ttype
	this.lattice = nil
}

func (this *Skeleton) Index() int {
	return this.index
}

func (this *Skeleton) Frame() int {
	return this.frame
}

func (this *Skeleton) SetIndex(index int) {
	this.index = index
}

func (this *Skeleton) Lattice(table *Table, index int) []int {
	if this.lattice == nil {
		this.lattice = make([][]int, 4)
		for _, ii := range []int{0, 1, 2, 3} {
			bb := table.block_config[this.ttype][ii]
			this.lattice[ii] = make([]int, len(bb))
			for i, v := range bb {
				if v != 0 {
					//不同玩家用不同颜色的方块
					this.lattice[ii][i] = v + (this.seatIndex + 1)
				} else {
					this.lattice[ii][i] = v
				}
			}
		}

	}
	return this.lattice[index]
}

func (this *Skeleton) Point() *Point {
	return this.point
}

func (this *Skeleton) Interval() *comp.Interval {
	return this.interval
}

func (this *Skeleton) Clone() *Skeleton {
	clone := &Skeleton{
		ttype:     this.ttype,
		index:     this.index,
		seatIndex: this.seatIndex,
		point: &Point{
			x: this.point.x,
			y: this.point.y,
		},
	}
	return clone
}

func (this *Skeleton) Move(x, y int) {
	this.point.x = x
	this.point.y = y
}

/**
检查方块是否发生碰撞
*/
func (this *Skeleton) CheckAABB(bb Skeleton) bool {
	return (this.point.x >= bb.Point().x && this.point.x < bb.Point().x+4 && this.point.y >= bb.Point().y && this.point.y < bb.Point().y+4)
}

func NewSkeletonQueue() *SkeletonQueue {
	this := &SkeletonQueue{
		SkeletonList: []*Skeleton{},
	}
	return this
}

type SkeletonQueue struct {
	SkeletonList []*Skeleton
}

func (this *SkeletonQueue) GetSkeleton(SeatIndex int) *Skeleton {
	for _, skeleton := range this.SkeletonList {
		if skeleton.GetSeatIndex() == SeatIndex {
			return skeleton
		}
	}
	return nil
}

func (this *SkeletonQueue) AddSkeleton(skeleton *Skeleton) error {
	this.SkeletonList = append(this.SkeletonList, skeleton)
	return nil
}

func (this *SkeletonQueue) RemoveSkeleton(SeatIndex int) error {
	for i, skeleton := range this.SkeletonList {
		if skeleton.GetSeatIndex() == SeatIndex {
			this.SkeletonList = append(this.SkeletonList[:i], this.SkeletonList[i+1:]...) // 最后面的“...”不能省略
		}
	}
	return nil
}

func (this *SkeletonQueue) ExecutionSkeleton(table *Table, player *Player, SeatIndex, Frame int, opcode string, msg map[string]interface{}) error {
	for _, skeleton := range this.SkeletonList {
		if skeleton.GetSeatIndex() == SeatIndex {
			//if (skeleton.Frame())>Frame{
			//	//这个操作是上一个方块的，直接丢弃
			//	log.Warning("ExecutionSkeleton timeout %v>%v",skeleton.Frame(),Frame)
			//	return	errors.Errorf("ExecutionSkeleton timeout %v>%v",skeleton.Frame(),Frame)
			//}
			if opcode == "RR" {
				table.grid.RotationRight(table, skeleton)
			} else if opcode == "MR" {
				table.grid.Move(table, skeleton, EAST)
			} else if opcode == "ML" {
				table.grid.Move(table, skeleton, WEST)
			} else if opcode == "MB" {
				//skeleton.Interval().SetStep(player.AccelerationSpeed())
				//快速向下
				var err error = nil
				for err == nil {
					err = table.grid.Move(table, skeleton, SOUTH)
				}
				//向上缓冲一行
				//table.grid.Move(table,skeleton,NORTH)
			} else if opcode == "MT" {
				//减速
				//skeleton.Interval().SetStep(player.GetSpeed())
			}
			skeleton.Start(table)
		}
	}
	return errors.New("events do not exist")
}

func (this *SkeletonQueue) Update(table *Table) error {
	kills := []int{}
	for _, skeleton := range this.SkeletonList {
		//player:=table.GetPlayerBySeatIndex(skeleton.GetSeatIndex())
		if skeleton.Interval().Step() {
			if table.grid.Move(table, skeleton, SOUTH) != nil {
				//向下移动失败,合并到方格中
				table.grid.Merge(table, skeleton)
				kills = append(kills, skeleton.GetSeatIndex())
			}
		}
	}
	for _, seatindex := range kills {
		this.RemoveSkeleton(seatindex)
	}
	return nil
}

type Block struct {
	value     int
	opacity   int
	operating int //控制码，备用
}

func (this *Block) Reset() {
	this.value = 0
	this.opacity = 0
	this.operating = 0
}
func (this *Block) Value() int {
	return this.value
}

func (this *Block) SetValue(Value int) {
	this.value = Value
}

func (this *Block) Opacity() int {
	return this.opacity
}

func (this *Block) SetOpacity(opacity int) {
	this.opacity = opacity
}

func (this *Block) Operating() int {
	return this.operating
}

func (this *Block) SetOperating(operating int) {
	this.operating = operating
}

// 游戏幕布
func NewGrid(width, height, winheight int) *Grid {
	this := &Grid{
		width:     width,
		height:    height,
		winheight: winheight,
		blocks:    make([]*Block, width*height),    //完整画布大小，在快速模式下画布高度会大于可见窗口高度
		render:    make([]*Block, width*winheight), //渲染画布
	}
	this.OnInitDataSync(this, 24)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if ((height - y) < 3) && x < width-3 {
				this.blocks[x+y*width] = &Block{}
			} else {
				this.blocks[x+y*width] = &Block{}
			}

		}
	}
	for x := 0; x < width; x++ {
		for y := 0; y < winheight; y++ {
			this.render[x+y*width] = &Block{}
		}
	}
	this.interval = comp.NewInterval(0, 500)
	this.slidingInterval = comp.NewInterval(0, 10000)
	return this
}

type Grid struct {
	component.DataSync
	width           int      //宽度
	height          int      //高度
	winheight       int      //窗口高度
	sliding         int      //滑动索引
	blocks          []*Block //块
	render          []*Block //画布块数据
	gameover        bool
	interval        *comp.Interval
	slidingInterval *comp.Interval
}

func (this *Grid) GameOver() bool {
	return this.gameover
}

func (this *Grid) SetSliding(sliding int) error {
	if (sliding + this.winheight) > this.height {
		return errors.New("Window limit")
	}
	if sliding < 0 {
		return errors.New("Window limit")
	}
	this.UpdateData()
	this.sliding = sliding
	return nil
}

func (this *Grid) MergeTemplate(s ...[][]int) (slice [][]int) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := this.MergeTemplate(s[1:]...) //...将数组元素打散
		slice = make([][]int, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}
	return
}

func (this *Grid) SetTemplate(table *Table, template [][]int) error {
	for y := len(template) - 1; y >= 0; y-- {
		//从最底层开始往上
		gridy := this.height - 1 - y
		trow := template[y]
		width := this.width
		if len(trow) < width {
			width = len(trow)
		}
		for x := 0; x < this.width; x++ {
			block := this.blocks[x+this.width*gridy]
			if x >= width {
				block.SetValue(1)
			} else {
				block.SetValue(trow[x])
			}

		}
	}
	return nil
}

func (this *Grid) Sliding() int {
	return this.sliding
}
func (this *Grid) GetBlock(table *Table, x, y int) *Block {
	if x < 0 || x >= this.width {
		//超出部分我们可以认为是被占了的
		return &Block{value: math.MaxInt64}
	}
	if y < 0 || y >= this.winheight {
		//超出部分我们可以认为是被占了的
		return &Block{value: math.MaxInt64}
	}
	index := x + ((y + this.sliding) * this.width)
	return this.blocks[index]
}

func (this *Grid) GetRender(table *Table, x, y int) *Block {
	if x < 0 || x >= this.width {
		//超出部分我们可以认为是被占了的
		return &Block{value: math.MaxInt64}
	}
	if y < 0 || y >= this.winheight {
		//超出部分我们可以认为是被占了的
		return &Block{value: math.MaxInt64}
	}
	index := x + y*this.width
	return this.render[index]
}

func (this *Grid) Move(table *Table, skeleton *Skeleton, direction [2]int) error {
	tox := skeleton.Point().X() + direction[0]
	toy := skeleton.Point().Y() + direction[1]
	if this.CheckAABB(table, skeleton.Lattice(table, skeleton.Index()), tox, toy) == false {
		return errors.Errorf("unable to move")
	}
	skeleton.Move(tox, toy)
	this.UpdateData()
	return nil
}

func (this *Grid) RotationRight(table *Table, skeleton *Skeleton) error {
	tox := skeleton.Point().X()
	toy := skeleton.Point().Y()
	index := (skeleton.Index() + 1) % 4
	if this.CheckAABB(table, skeleton.Lattice(table, index), tox, toy) == false {
		return errors.Errorf("unable to move")
	}
	skeleton.SetIndex(index)
	this.UpdateData()
	return nil
}

func (this *Grid) Update(table *Table) error {
	if this.slidingInterval.Step() {
		//往上滚动
		this.SetSliding(this.Sliding() + 1)
	}

	return nil
}

/*
下沉一行
*/
func (this *Grid) Sink(table *Table, index int) error {
	for y := this.winheight - 1; y >= 0; y-- {
		//从最底层开始检查是否正行被填充完整了
		if y <= index {
			for x := 0; x < this.width; x++ {
				block := this.GetBlock(table, x, y)
				prerowblock := this.GetBlock(table, x, y-1)
				if prerowblock.Value() == math.MaxInt64 {
					//超出界线部分
					block.SetValue(0)
				} else {
					block.SetValue(prerowblock.Value())
				}
			}
		}
	}
	return nil
}

/**
检查方块是否发生碰撞
*/
func (this *Grid) CheckAABB(table *Table, lattice []int, startx, starty int) bool {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			index := x + y*4
			if lattice[index] != 0 {
				block := this.GetBlock(table, x+startx, y+starty)
				if block.Value() != 0 {
					//有冲突
					return false
				}
			}
		}
	}
	return true
}

/**
合并到方格中
*/
func (this *Grid) Merge(table *Table, skeleton *Skeleton) error {
	this.UpdateData()
	//将骨骼中有方块的块覆盖画布
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			index := x + y*4
			block := this.GetBlock(table, x+skeleton.Point().X(), y+skeleton.Point().Y())
			if skeleton.Lattice(table, skeleton.Index())[index] > 0 {
				block.SetValue(skeleton.Lattice(table, skeleton.Index())[index])
			}
		}
	}

	for y := 0; y < this.winheight; y++ {
		//从最底层开始检查是否正行被填充完整了
		complete := true
		for x := 0; x < this.width; x++ {
			block := this.GetBlock(table, x, y)
			if block.Value() == 0 {
				complete = false
				break
			}
		}
		if complete {
			player := table.GetPlayerBySeatIndex(skeleton.GetSeatIndex())
			player.SetScore(player.GetScore() + 15)
			this.Sink(table, y)
			table.S2C_GameStatusChangeBroadcast()
		}
	}

	// 检查游戏是否已结束
	for x := 0; x < this.width; x++ {
		block := this.GetBlock(table, x, 0)
		if block.Value() != 0 {
			this.gameover = true
			break
		}
	}

	return nil
}

/**
合并渲染
*/
func (this *Grid) MergeRender(table *Table) error {
	//先清零 合并底层画布
	for y := 0; y < this.winheight; y++ {
		for x := 0; x < this.width; x++ {
			block := this.GetRender(table, x, y)
			bb := this.GetBlock(table, x, y)
			block.SetValue(bb.Value())
			if bb.Value() == 0 {
				block.SetOpacity(0)
			} else {
				block.SetOpacity(255)
			}
		}
	}

	//合并移动中的方块
	for _, skeleton := range table.skeletonQueue.SkeletonList {
		//投影计算
		//clone := skeleton.Clone()
		//var err error = nil
		//for err == nil {
		//	err = this.Move(table, clone, SOUTH)
		//}
		//if skeleton.Point().Y()+4<clone.Point().Y(){
		//相互没有交集才绘制投影
		//for x:=0;x<4;x++{
		//	for y:=0;y<4;y++{
		//		index:=x+y*4
		//		block:=this.GetRender(table,x+clone.Point().X(),y+clone.Point().Y())
		//		if block.Value()==0{
		//			v:=clone.Lattice(table,clone.Index())[index]
		//			if v>0{
		//				block.SetValue(v)
		//				block.SetOpacity(30)
		//			}
		//		}
		//	}
		//}
		//}
		for x := 0; x < 4; x++ {
			for y := 0; y < 4; y++ {
				index := x + y*4
				block := this.GetRender(table, x+skeleton.Point().X(), y+skeleton.Point().Y())
				v := skeleton.Lattice(table, skeleton.Index())[index]
				if v > 0 {
					block.SetValue(v)
					block.SetOpacity(255)
				}
			}
		}
	}

	return nil
}

func (this *Grid) Progress(table *Table) float64 {
	fm := float64(this.height - this.winheight)
	if fm > 0.0 {
		return float64(this.sliding) / fm
	} else {
		return 100.0
	}
}

func (this *Grid) Source(table interface{}) ([]byte, error) {
	this.MergeRender(table.(*Table))
	render := &pb.S2C_GridBroadcast{}
	render.Frame = *proto.Int64(int64(1))
	render.Width = *proto.Int64(int64(this.width))
	render.Height = *proto.Int64(int64(this.winheight))
	render.Map = make([]*pb.Block, this.width*this.height)
	for i, block := range this.render {
		render.GetMap()[int32(i)] = &pb.Block{
			// 使用辅助函数设置域的值
			Value:     *proto.Uint32(uint32(block.Value())),
			Opacity:   *proto.Uint32(uint32(block.Opacity())),
			Operating: *proto.Uint32(uint32(block.Operating())),
			Index:     *proto.Uint32(uint32(i)),
		}
	}
	return proto.Marshal(render)
}
