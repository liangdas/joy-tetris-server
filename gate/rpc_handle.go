package mgate

import (
	"fmt"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/selector"
	"net/url"
)

// Selector 客户端路由规则自定义函数
func (this *Gate) Selector(session gate.Session, topic string, u *url.URL) (s module.ServerSession, err error) {
	moduleType := u.Scheme
	nodeId := u.Hostname()
	//使用自己的
	if nodeId == "modulus" {
		//取模
	} else if nodeId == "cache" {
		//缓存
	} else if nodeId == "random" {
		//随机
	} else {
		//
		//指定节点规则就是 module://[user:pass@]nodeId/path
		//方式1
		//moduleType=fmt.Sprintf("%v@%v",moduleType,u.Hostname())
		//方式2
		serverID := fmt.Sprintf("%v@%v", moduleType, nodeId)
		return this.GetRouteServer(moduleType, selector.WithFilter(func(services []*registry.Service) []*registry.Service {
			for _, service := range services {
				for _, node := range service.Nodes {
					if node.Id == serverID {
						return []*registry.Service{service}
					}
				}
			}
			return []*registry.Service{}
		}))
	}
	return this.GetRouteServer(moduleType)
}
