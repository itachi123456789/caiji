package main

import (
	"./common"
	"./worker"
	"log"
	"time"
)

func main() {
	var starton bool = true //自动识别对方网站是否维护 并根据维护时间关闭或开启 采集
	langx := "zh-cn"
	/*
		2项都为FS时为冠军赛事
		FU足球早餐,FT为今日足球，rtype分类：pd是全场波胆，hpd是上半场胆球，r是独赢、让球、大小、单双，t是总入球，f是半场、全场
		BK是今日 篮球/美式足球， BKR篮球早餐/美式足球 ，rtype分类：r_main独赢、让球、大小
		足球滚球,FT rtype分类：re是独赢、让球、大小、单双
		BK，篮球滚球，re_main独赢、让球、大小
		TN网球，r_all 独赢、让盘、大小
		VB排球，r_all 独赢、让盘、大小
		BM羽毛球，r_all 独赢、让盘、大小
		OP其他，r 独赢、让盘、大小、单双
		gtype:SAIGUO ,rtype:FT 为足球赛果,BK为篮球赛果,TN为网球,VB为排球，BM为羽毛球
	*/
	go func() {
		for {
			select {
			case <-common.SendChannel:
				go common.Sendmsg()
			case <-common.LogsChannel:
				go common.Logs()
			case <-common.ErrlogsChannel:
				go common.LogErr()
			}
		}
	}()
	common.SendChannel <- true
	common.LogsChannel <- true
	common.ErrlogsChannel <- true
	/*
		gtype := "FU"
		rtype := "r"
		for {
			go worker.Theif_data(gtype, rtype, langx, "0")
			time.Sleep(1e9 * 50)
		}

	*/
	var gtype [10]string
	var rtype [10][6]string

	gtype[0] = "FU"
	rtype[0][0] = "r"
	rtype[0][1] = "pd"
	rtype[0][2] = "hpd"
	rtype[0][3] = "t"
	rtype[0][4] = "f"

	gtype[1] = "FT"
	rtype[1][0] = "r"
	rtype[1][1] = "pd"
	rtype[1][2] = "hpd"
	rtype[1][3] = "t"
	rtype[1][4] = "f"

	gtype[2] = "BK"
	rtype[2][0] = "r_main"

	gtype[3] = "BKR"
	rtype[3][0] = "r_main"

	gtype[4] = "BM"
	rtype[4][0] = "r_all"

	gtype[5] = "VB"
	rtype[5][0] = "r_all"

	gtype[6] = "TN"
	rtype[6][0] = "r_all"

	gtype[7] = "SAIGUO"
	rtype[7][0] = "FT"
	rtype[7][1] = "BK"
	rtype[7][2] = "TN"
	rtype[7][3] = "VB"
	rtype[7][4] = "BM"
	rtype[7][5] = "OP"

	gtype[8] = "FS"
	rtype[8][0] = "FS"

	gtype[9] = "OP"
	rtype[9][0] = "r"
	worker.Getjieshuiurl(langx) //获取接水可能出现的新网址，这个页面会得到一个cookie，但是现时这个cookie可有可无，本程序已带上这个cookie
	log.Println("程序正在运行......")

	for {
		select {
		case <-worker.Other_time.C:
			if starton {
				for i := 0; i < len(gtype); i++ {
					for ii := 0; ii < len(rtype[i]) && rtype[i][ii] != ""; ii++ {
						worker.Theif_data(gtype[i], rtype[i][ii], langx, "0") //采集入口
					}

				}
			}
			worker.Other_time.Reset(time.Duration(worker.Other))
		case <-worker.Zq_time.C: //足球滚球每秒采集一次
			if starton {
				go worker.Theif_data("FT", "re", langx, "0")
			}
			worker.Zq_time.Reset(time.Duration(worker.Zqgq))
		case <-worker.Lq_time.C: //篮球滚球每秒采集一次
			if starton {
				go worker.Theif_data("BK", "re_main", langx, "0")
			}
			worker.Lq_time.Reset(time.Duration(worker.Lqgq))
		case <-worker.Time10s.C: //十秒检查一次
			if common.Logtimes < 10 { //tcp发送协程重启的日志次数 防止网络问题导致不断写入日志 每十秒增加一个日志写入机会
				common.Logtimes = common.Logtimes + 1
			}
			if worker.MaxAgain < 10 && time.Now().Minute() == 18 { //每小时增加 uid可重新获取的次数
				log.Println(worker.MaxAgain)
				worker.MaxAgain = worker.MaxAgain + 1
			}

			worker.Readcf("config.ini")                   //读配置文件
			if worker.Zantingshijian.Before(time.Now()) { //麻痹的皇冠规则变了 原来在未登录前的页面显示维护以及维护时间的 现在在内页 靠 这个功能先不做 留着不影响
				starton = true
				worker.Zq_time.Reset(time.Second)
				worker.Lq_time.Reset(time.Second)
			} else {
				starton = false
				worker.Zq_time.Stop()
				worker.Lq_time.Stop()
			}
			worker.Time10s.Reset(time.Second * 10)
		}
	}

}
