package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
)

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Get() {
	c.Response.Writeln("Controller Show")

}

func InitControl() {
	s := g.Server()
	s.SetServerRoot(".")
	s.SetRewrite("/index", "/index.html")
	s.SetRewrite("/", "/index.html")

	ctl := new(Controller)
	g2 := s.Group("/api")
	g2.REST("/handler", ctl)

	s.SetPort(1011)
	s.Run()
}
