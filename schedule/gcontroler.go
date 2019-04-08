package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/frame/gmvc"
	"strconv"
)

func init() {

	g.Config().SetFileName("config.json")

	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		Host:    g.Config().GetString("database.default.0.host"),
		Port:    g.Config().GetString("database.default.0.port"),
		User:    g.Config().GetString("database.default.0.user"),
		Pass:    g.Config().GetString("database.default.0.pass"),
		Name:    g.Config().GetString("database.default.0.name"),
		Type:    g.Config().GetString("database.default.0.type"),
		Role:    "master",
		Charset: "utf8",
	})
	var err error
	db, err = gdb.New()
	if err != nil {
		panic(err)
	}
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)
}

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Get() {

	page := c.Request.Get("page")
	fmt.Println(page)
	i, _ := strconv.Atoi(page)
	if i < 1 {
		i = 1
	}

	r, _ := db.Table("sys_app_function_time").Cache(3, "sys_app_function_time"+page).Where("function_id like 'HWW%'").ForPage(i, 10).OrderBy("id desc").Select()
	fmt.Println(r.ToJson())
	c.Response.Writeln(r.ToJson())

}

func InitControl() {
	s := g.Server()
	s.SetServerRoot(".")
	s.SetRewrite("/index", "/index.html")
	s.SetRewrite("/", "/index.html")

	ctl := new(Controller)
	g2 := s.Group("/api")
	g2.REST("/handler/:page", ctl)

	s.SetPort(1011)
	s.Run()
}
