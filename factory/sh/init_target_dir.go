package sh

import (
	"ek8/common"
	"ek8/config"
	"ek8/shell"
)

func Init_target_dir(allServers map[string]string) {
	if common.ModFlags["config"] != 1 {
		return
	}
	common.TPrintln("init target [config]...")
	for _, ip := range allServers {
		target_dir := config.Cfg["release_dir"] + "/" + ip

		cmd := new(common.Command).
			Add("source config.cfg").
			//Add("rm -rf " + target_dir + "/cfg").
			Add("mkdir -p " + target_dir + "/cfg/base").
			Add("\\cp -f config.cfg " + target_dir).
			Add("\\cp -f $sub_conf " + target_dir + "/cfg").
			Add("\\cp -f cfg/base/*.config " + target_dir + "/cfg/base/").
			Add("echo 完成").
			ToString()
		state, _ := shell.Exec(cmd, ip)
		if !state {
			panic(ip + ": 初始化目录" + "出错！")
		}
	}
}
