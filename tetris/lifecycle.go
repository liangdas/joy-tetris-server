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
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"time"
)

func (this *Table) OnCreate() {
	this.QTable.OnCreate()
	this.current_frame = 0
	this.sync_frame = 0
}
func (this *Table) OnDestroy() {
	this.QTable.OnDestroy()
}

func (mg *Table) GetApp() module.App {
	return mg.module.GetApp()
}

//在table超时是调用
func (this *Table) OnTimeOut() {
	this.Finish()
}

func (table *Table) GetSeats() map[string]room.BasePlayer {
	m := map[string]room.BasePlayer{}
	for _, v := range table.players {
		m[fmt.Sprintf("%v", v.GetSeatIndex())] = v
	}
	return m
}

func (this *Table) GetModule() module.App {
	return this.module.GetApp()
}

// Update 定帧计算所有玩家的位置
func (this *Table) Update(ds time.Duration) {
	err := this.fsmEvent(GoForward, this)
	if err != nil {
		log.Error("trigger err: %v", err)
	}
}
