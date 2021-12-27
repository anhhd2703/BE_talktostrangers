package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	sig "os/signal"
	"syscall"
	server "talktostrangers/cmd/biz/json-rpc"

	log "github.com/pion/ion-log"
	"github.com/spf13/viper"
)

var (
	conf server.SignalServer
	file string
)

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -h (show help info)")
}

func unmarshal(rawVal interface{}) bool {
	if err := viper.Unmarshal(rawVal); err != nil {
		fmt.Printf("config file %s loaded failed. %v\n", file, err)
		return false
	}
	return true
}

func load() bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file %s read failed. %v\n", file, err)
		return false
	}

	if !unmarshal(&conf) || !unmarshal(&conf) {
		return false
	}
	if err != nil {
		fmt.Printf("config file %s loaded failed. %v\n", file, err)
		return false
	}

	fmt.Printf("config %s load ok!\n", file)

	return true
}

func parse() bool {
	flag.StringVar(&file, "c", "conf/conf.toml", "config file")
	help := flag.Bool("h", false, "help info")
	flag.Parse()
	if !load() {
		return false
	}

	if *help {
		showHelp()
		return false
	}
	return true
}

func main() {
	if !parse() {
		showHelp()
		os.Exit(-1)
	}
	s := server.NewServer(conf.Config)
	if err := s.Start(); err != nil {
		log.Errorf("biz start error: %v", err)
		os.Exit(-1)
	}
	// defer s.Close()

	// Press Ctrl+C to exit the process
	ch := make(chan os.Signal, 1)
	sig.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
