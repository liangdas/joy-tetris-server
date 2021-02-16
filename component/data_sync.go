package component

import (
	"github.com/liangdas/mqant/log"
	"hash/crc32"
	"joysrv/component/rsync"
	"time"
)

const (
	FULL = iota
	PATCH
	NEWEST //无需同步
)

type InterDataSync interface {
	ResetData() error
	UpdateData()
	SyncDate()
	Marshal(table interface{}) ([]byte, int, error)
}

type SyncBytes interface {
	//Marshal() ([]byte,int, error)		//编辑数据
	Source(table interface{}) ([]byte, error) //数据源
}

type DataSync struct {
	sub        SyncBytes
	original   []byte //上一次同步数据
	syncDate   int64  //数据同步给客户端的日期
	updateData int64  //数据更新日期
}

func (this *DataSync) OnInitDataSync(Sub SyncBytes) error {
	this.sub = Sub
	return nil
}

/**
重置补丁
*/
func (this *DataSync) ResetData() error {
	this.original = nil
	this.updateData = time.Now().UnixNano()
	return nil
}
func (this *DataSync) UpdateData() {
	this.updateData = time.Now().UnixNano()
}
func (this *DataSync) SyncDate() {
	this.syncDate = time.Now().UnixNano()
}

/**
补丁数据
*/
func (this *DataSync) Marshal(table interface{}) ([]byte, int, error) {
	if this.updateData <= this.syncDate {
		return nil, NEWEST, nil
	}
	modified, err := this.sub.Source(table)
	if err != nil {
		return nil, 0, err
	}
	if this.original == nil {
		this.original = modified
		this.SyncDate()
		log.TInfo(nil, "this.original=modified")
		return modified, FULL, err
	}
	rs := &rsync.LRsync{
		BlockSize: 24,
	}
	hashes := rs.CalculateBlockHashes(this.original)
	opsChannel := rs.CalculateDifferences(modified, hashes)
	ieee := crc32.NewIEEE()
	ieee.Write(modified)
	s := ieee.Sum32()
	delta := rs.CreateDelta(opsChannel, len(modified), s)
	this.original = modified
	this.SyncDate()
	return delta, PATCH, err
}
