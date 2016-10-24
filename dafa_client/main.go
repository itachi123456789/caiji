package main

import (
	"./common"
	"./worker"
	"github.com/kardianos/service"
	"log"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	log.Println("Dafa Game History is Running!")
	go common.LogErr()   //写错误日志协程
	go common.Logs()     //写正常日志协程
	go common.PostData() //post提交回数据库服务器
	//go worker.KgStart()  //kg彩票协程
	//go worker.AgStart() //ag真人游戏协程
	go worker.SpStart() //sp体育协程

	for {
		select {
		case <-common.PostData_watch: //监视提交入库协程PostData()
			go common.PostData()

		case <-worker.KgStart_watch: //99%用不到的 监视KgStart()协程
			go worker.KgStart()

		case <-worker.SpStart_watch: //99%用不到的 监视SpStart()协程
			go worker.SpStart()

		case <-worker.AgStart_watch: //99.99%用不到的 监视AgStart()协程
			go worker.AgStart()

		case <-common.ErrLogs_watch: //98%用不到的 监视写错误日志LogErr()协程
			go common.LogErr()

		case <-common.Logs_watch: //98%用不到的 监视写正常日志Logs()协程
			go common.Logs()
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	prg := &program{}
	svcConfig := &service.Config{
		Name:        "dafaGameHistory",
		DisplayName: "Dafa Game History",
		Description: "采集大发游戏历史记录的客户端",
	}

	s, err := service.New(prg, svcConfig)

	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
