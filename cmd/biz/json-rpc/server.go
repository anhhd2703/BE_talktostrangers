package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"talktostrangers/pkg/biz"
	db "talktostrangers/pkg/db/model"
	"time"
)

type SignalServer struct {
	Config    signalConf `mapstructure:"signal"`
	redis     *db.Redis
	Router    *http.ServeMux
	WareHouse *biz.WareHouse
}

type signalConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	SignalConnect
}
type SignalConnect struct {
	Redis db.Config `mapstructure:"redis"`
	Etcd  etcdConf  `mapstructure:"etcd"`
	Log   logConf   `mapstructure:"log"`
}
type logConf struct {
	Level string `mapstructure:"level"`
}
type etcdConf struct {
	Addrs []string `mapstructure:"addrs"`
}

func NewServer(config signalConf) *SignalServer {
	return &SignalServer{
		Config:    config,
		WareHouse: biz.NewWareHouse(),
	}
}

func (s *SignalServer) startRedis() {
	log.Println("data", s.Config)
	s.redis = db.NewRedis(s.Config.SignalConnect.Redis)
}

func (s *SignalServer) Start() error {
	s.startRedis()
	// use gorilla mux to middle ware
	//Use the default DefaultServeMux.
	fmt.Println(s.Config)
	s.Router = http.NewServeMux()
	s.initializeRoutes()
	myHttp := &http.Server{
		Addr:           s.Config.Host + ":" + strconv.Itoa(s.Config.Port),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.Router,
	}
	err := myHttp.ListenAndServe()
	return err
}

func (s *SignalServer) initializeRoutes() {
	s.Router.Handle("/room", biz.GetRooomHandler(s.WareHouse))
}
