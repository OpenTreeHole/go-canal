package main

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/sirupsen/logrus"
	"os"
)

var log = &logrus.Logger{
	Out: os.Stdout,
	Formatter: &logrus.TextFormatter{
		DisableQuote: true,
	},
	Hooks: make(logrus.LevelHooks),
	Level: logrus.DebugLevel,
}

func main() {
	go StartTimer()

	c, err := canal.NewCanal(NewConfig())
	if err != nil {
		panic(err)
	}

	c.SetEventHandler(&MyEventHandler{})

	position, err := c.GetMasterPos()
	if err != nil {
		panic(err)
	}

	err = c.RunFrom(position)
	if err != nil {
		panic(err)
	}
}
