/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package webapp

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

// Module web 模块
var Module = func() *WebApp {
	web := new(WebApp)
	return web
}

// Web 结构对象基于 BaseModule
type WebApp struct {
	basemodule.BaseModule
	StaticPath   string
	Port         int
}

// GetType 获取模块类型标识
func (self *WebApp) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "webapp"
}

//Version 获取Web模块版本号
func (self *WebApp) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

//OnInit Web模块初始化方法
func (self *WebApp) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.StaticPath = self.GetModuleSettings().Settings["StaticPath"].(string)
	self.Port = int(self.GetModuleSettings().Settings["Port"].(float64))
}


func registerFilter(e *echo.Echo) {
	// middleware
	e.Use(middleware.Recover())
}

//Run Web模块启动方法
func (self *WebApp) Run(closeSig chan bool) {
	//这里如果出现异常请检查8080端口是否已经被占用

	e := echo.New()
	registerFilter(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
	e.Static("/static", "static")
	go func() {
		log.Info("webapp server Listen : %s", fmt.Sprintf(":%d", self.Port))
		// Start server
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", self.Port)))
	}()

	<-closeSig
	log.Info("webapp server Shutting down...")
	e.Close()
}

//OnDestroy Web模块注销方法
func (self *WebApp) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}
