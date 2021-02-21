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
package comp

import (
	"time"
)

type IntervalInterface interface {
	IntervalInit(Milliseconds, Step int64)
	Reset()
	Rest() time.Duration
	Readied() bool  //准备好的
	Step() bool     //步进一次
	Complete() bool //是否完成
	End() bool      //是否结束
}

/**
定时间隔
*/
type Interval struct {
	Milliseconds int64 //多少毫米后结束  如果0则无限循环
	step         int64 //每个多少毫米触发一次步调
	preStep      int64 //上一次步调时间
	start        int64 //开始时间
	end          bool
}

func NewInterval(Milliseconds, Step int64) *Interval {
	this := &Interval{
		Milliseconds: Milliseconds,
		step:         Step,
		end:          false,
	}
	this.start = time.Now().UnixNano() / 1e6
	return this
}

func (this *Interval) IntervalInit(Milliseconds, Step int64) {
	this.Milliseconds = Milliseconds
	this.step = Step
	this.end = false
	this.start = time.Now().UnixNano() / 1e6
}

/**
重置，当定时器结束后，想重新启用则需要重置
*/
func (this *Interval) Reset() {
	this.start = time.Now().UnixNano() / 1e6
	this.preStep = -1
	this.end = false
}

/**
还剩余多少时间结束
*/
func (this *Interval) Rest() time.Duration {
	now := time.Now().UnixNano() / 1e6
	rest := this.start + this.Milliseconds - now
	if rest <= 0 {
		return time.Duration(0)
	}
	return time.Millisecond * time.Duration(rest)
}

/**
检查是否已经达到下一个步骤[不会更新步骤时间状态]
*/
func (this *Interval) Readied() bool {
	now := time.Now().UnixNano() / 1e6
	if (now - this.preStep) > this.step {
		return true
	}
	return false
}

/**
是否已经达到下一个步骤[步骤状态将被设置为最新时间，下一次调用Step()为true将是在下一个步骤时间到达以后]
*/
func (this *Interval) Step() bool {
	if this.end {
		return false
	}
	now := time.Now().UnixNano() / 1e6
	if (now - this.preStep) > this.step {
		this.preStep = now
		return true
	}
	return false
}

/**
当定时器达到完成条件返回true, 该函数在时间到达后只会返回一次true，如果要重新开始定时需要执行Reset()方法
*/
func (this *Interval) Complete() bool {
	now := time.Now().UnixNano() / 1e6
	if (now - this.start) > this.Milliseconds {
		this.end = true
		return true
	}
	return false
}
func (this *Interval) End() bool {
	return this.end
}
