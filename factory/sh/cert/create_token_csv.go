package cert

import (
	"ek8/common"
	"ek8/shell"
)

func Create_token_csv(allservers map[string]string) {
	if common.ModFlags["tokenCSV"] != 1 {
		return
	}
	common.TPrintln("创建token_csv...")

	//创建token.csv文件
	token_file := "target/.token.csv"
	cmd := new(common.Command).
		Add("source config.cfg").
		Add("export bootstrap_token=\"$(head -c 6 /dev/urandom | md5sum |head -c 6).$(head -c 16 /dev/urandom | md5sum | head -c 16)\"").
		Add("echo \"$bootstrap_token,${bootstrap_user},10001,\\\"${bootstrap_group}\\\"\" > " + token_file).
		ToString()
	shell.Exec(cmd, "local")

	for _, ip := range allservers {
		allca_dir := "target/" + ip + "/cert"
		kube_bootstrap_dir_name := "bootstrap"
		kube_bootstrap_dir := allca_dir + "/" + kube_bootstrap_dir_name

		//拷贝token文件到 target/ip目录
		cmd := new(common.Command).
			Add("source config.cfg").
			Add("mkdir -p " + kube_bootstrap_dir).
			Add("cp -f " + token_file + " " + kube_bootstrap_dir + "/token.csv").
			Add("echo \"token(svc) file 创建.....\"").
			ToString()
		shell.Exec(cmd, ip)

	}
	//清理临时文件
	//cmd = new(common.Command).
	//	Add("rm -rf " + token_file).ToString()
	//shell.Exec(cmd, "local")

}
