package config

import "flag"

type sFlag struct {
	Help bool
	ConfigFile string
}

var Flag sFlag

func ParseFlag() {
	flag.BoolVar(&Flag.Help, "h", false, "打印帮助")
	flag.StringVar(&Flag.ConfigFile, "c", "./config.toml", "配置文件路径")
	flag.Parse()
}
