package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"sync"
)

func Create_sa(allservers map[string]string) {
	if common.ModFlags["saCERT"] != 1 {
		return
	}
	common.TPrintln("创建 sa...")
	var once sync.Once
	for _, ip := range allservers {
		allca_dir := "$PWD/target/" + ip + "/cert"
		sa_dir_name := "sa"
		sa_dir := allca_dir + "/" + sa_dir_name

		sa_key_file := "sa.key"
		sa_pub_file := "sa.pub"

		//创建sa (service account)
		createCert := func() {
			cmd := new(common.Command).Add("source config.cfg").
				Add("mkdir  -p " + sa_dir).
				Add("mkdir -p target/.cert/" + sa_dir_name).
				Add("cd target/.cert/" + sa_dir_name).
				Add("openssl genrsa -out " + sa_key_file + " 2048").
				Add("openssl rsa -in " + sa_key_file + " -pubout -out " + sa_pub_file).
				Add("cd $PWD").
				Add("echo \"sa 创建.....\"").
				ToString()
			shell.Exec(cmd, "local")
		}
		once.Do(createCert)
		assist.Copy_cert_from_local_to_target("target/.cert", sa_dir_name, "sa", sa_dir_name, ip)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  target/.cert")
	}
}
