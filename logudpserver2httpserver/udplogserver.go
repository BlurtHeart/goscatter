package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/donnie4w/go-logger/logger"
	"github.com/mreiferson/go-options"
)

// 创建通道
var c chan []byte = make(chan []byte, 100000)

// 收到的包的总数
var TotalRecvCount int64

var (
	flagSet = flag.NewFlagSet("Active6 Glj Transfer", flag.ExitOnError)
	config  = flagSet.String("config", "", "path to config file")

	listenIP   = flagSet.String("listen-ip", "127.0.0.1", "listen ip")
	listenPort = flagSet.Int("listen-port", 8000, "listen port")

	console     = flagSet.Bool("console", true, "print log output console")
	level       = flagSet.Int("level", 0, "print log by set level. 0:default,1:debug,2:info,3:warn,4:error,5:fatal")
	logdir      = flagSet.String("logdir", "", "save log file in logdir")
	logfilename = flagSet.String("logfilename", "", "print log write to logfilename")
	lognum      = flagSet.Int("lognum", 5, "number of logs")
	logfilesize = flagSet.Int64("logfilesize", 10, "size of logfile")

	hosturl  = flagSet.String("hosturl", "", "figure out the remote addr data will send to")
	topic    = flagSet.String("topic", "", "database data will store")
	user     = flagSet.String("user", "", "authentication user of the remote server")
	password = flagSet.String("password", "", "password of user")
)

type Options struct {
	ListenIP   string `flag:"listen-ip"`   // udp服务监听ip
	ListenPort int    `flag:"listen-port"` // udp服务监听端口

	Console     bool   `flag:"console"`     // 是否打印到控制台
	Level       int    `flag:"level"`       // 日志等级：ALL、INFO、DEBUG、WARN、ERROR、FATAL、OFF
	LogDir      string `flag:"logdir"`      // 存放文件日志目录
	LogFileName string `flag:"logfilename"` // 保存日志的文件名
	LogNum      int32  `flag:"lognum"`      // 日志文件数
	LogFileSize int64  `flag:"logfilesize"` // 每个日志文件大小

	Hosturl  string `flag:"hosturl"`  // 远端服务器的地址
	Topic    string `flag:"topic"`    // 数据表名
	User     string `flag:"user"`     // 远端服务器的认证用户
	Password string `flag:"password"` // 密码
}

func Stated() {
	var nodesend int64

	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			totalCount := atomic.LoadInt64(&TotalRecvCount)
			sendOneCount := totalCount - nodesend
			logger.Info("Per Second Static: TotalRecvCounts:(" + strconv.FormatInt(totalCount, 10) + ")\t 1Sec RecvCounts:(" + strconv.FormatInt(sendOneCount, 10) + ")")
			nodesend = totalCount
		}
	}
}

func main() {

	//读取配置
	flagSet.Parse(os.Args[1:])
	var cfg map[string]interface{}
	if *config != "" {
		_, err := toml.DecodeFile(*config, &cfg)
		if err != nil {
			log.Fatalf("ERROR: failed to load config file %s - %s", *config, err.Error())
		}
	}
	opts := &Options{}
	options.Resolve(opts, flagSet, cfg)

	// 创建Logger
	logger.SetConsole(true)
	logger.SetRollingFile(opts.LogDir, opts.LogFileName, opts.LogNum, opts.LogFileSize, logger.MB)
	logger.SetLevel(logger.INFO)

	logger.Warn("begin ...")

	// 创建监听
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.ParseIP(opts.ListenIP),
		Port: opts.ListenPort,
	})
	if err != nil {
		fmt.Println("监听失败!", err)
		return
	}
	defer socket.Close()

	go func() {
		Stated()
	}()

	go httpPost(opts)

	for {
		// 读取数据
		data := make([]byte, 4096)
		dataCount, _, err := socket.ReadFromUDP(data)
		//		buff := make([]byte, dataCount)
		//		copy(buff[:], data)
		if err != nil {
			fmt.Println("读取数据失败!", err)
			continue
		}

		// 将数据放入通道
		c <- data[:dataCount]

		atomic.AddInt64(&TotalRecvCount, 1)
	}

	logger.Warn("finish ...")
}

func httpPost(opts *Options) {
	client := http.Client{}
	for {
		data := <-c
		req, _ := http.NewRequest("POST", opts.Hosturl, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Add("User", opts.User)
		req.Header.Add("Password", opts.Password)
		req.Header.Add("Topic", opts.Topic)
		req.Header.Add("Row-Split", "\n")
		req.Header.Add("Field-Split", ",")
		req.Header.Add("Format", "csv")
		resp, err := client.Do(req)
		defer resp.Body.Close()
		response, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(response), err)
	}

}
