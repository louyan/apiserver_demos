package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	v "yannotes.cn/apiserver_demos/demo10/pkg/version"
	"yannotes.cn/apiserver_demos/demo10/router/middleware"

	"yannotes.cn/apiserver_demos/demo10/model"

	"github.com/lexkong/log"

	"net/http"
	"time"

	"github.com/spf13/viper"

	"yannotes.cn/apiserver_demos/demo10/config"

	"github.com/spf13/pflag"

	"github.com/gin-gonic/gin"
	"yannotes.cn/apiserver_demos/demo10/router"
)

/*
每秒查询数 QPS 是对一个特定的查询服务器在规定时间内所处理流量多少的衡量标准
QPS = 并发数 / 平均请求响应时间*/
var (
	cfg     = pflag.StringP("config", "c", "", "apiserver config file path.")
	version = pflag.BoolP("version", "v", false, "sjhow version info")
)

func main() {
	pflag.Parse()
	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return

	}

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	/*for {
		//热更新测试
		//fmt.Println(viper.GetString("runmode"))
		//time.Sleep(4 * time.Second)
		//日志测试
		log.Info("1111111111111111111111111")
		time.Sleep(100 * time.Millisecond)
	}*/

	// Set gin mode
	gin.SetMode(viper.GetString("runmode"))

	// Create the Gin engine
	g := gin.New()

	// gin middlewares
	//middlewares := []gin.HandlerFunc{}

	router.Load(
		// Cores.
		g,
		// Middlewares

		middleware.Logging(),
		middleware.RequestId(),
		//middlewares...,
	)

	// Ping the server to make sure the router is working
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	//测试日志形式
	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	//log.Warnf("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	//log.Debugf("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	//log.Errorf(errors.New("error"), "Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
	//log.Fatalf(errors.New("error"), "Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}
