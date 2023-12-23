package service

import (
	"ek8/common"
	"ek8/config"
)

func Create_coredns_conf(apiServers map[string]string) {
	if common.ModFlags["corednsCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建coredns配置...")

	var ip string
	for _, v := range common.MapToArray(apiServers) {
		ip = v
		break
	}

	target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/coredns"

	coredns_yaml_sed_file_name := "coredns.yaml.sed"
	deploy_sh_file_name := "deploy.sh"
	dnsutils_yaml_file_name := "dnsutils.yaml"
	coredns_yaml_sed_file := target_dir + "/" + coredns_yaml_sed_file_name
	deploy_sh_file := target_dir + "/" + deploy_sh_file_name
	dnsutils_yaml_file := target_dir + "/" + dnsutils_yaml_file_name

	template_coredns_yaml_sed_file := "factory/template/service/coredns/" + coredns_yaml_sed_file_name
	template_deploy_sh_file := "factory/template/service/coredns/" + deploy_sh_file_name
	template_dnsutils_yaml_file := "factory/template/service/coredns/" + dnsutils_yaml_file_name

	//cmd := new(common.Command).
	//	Add("source config.cfg").
	//	Add("rm -rf " + target_dir).
	//	Add("mkdir -p " + target_dir).
	//	Add("\\cp -f factory/template/coredns/coredns.yaml.sed " + current_conf_file).
	//	Add("\\cp -f  factory/template/coredns/deploy.sh " + current_sh_file).
	//	Add("\\cp -f  factory/template/coredns/dnsutils.yaml " + target_dir + "/").
	//	Add("sed -i " + config.SED_CHAR + " \"s/{{COREDNS_VERSION}}/${coredns_version//\\//\\\\/}/g\" " + current_conf_file).
	//	Add("echo " + ip + " config 完成").
	//	ToString()
	//state, _ := shell.Exec(cmd, "ip")
	//if !state {
	//	panic(ip + ": 创建coredns配置出错！")
	//}

	//创建 coredns.yaml.sed
	new(common.FileAssist).
		Load(template_coredns_yaml_sed_file).
		Println("创建 "+coredns_yaml_sed_file_name+"...", ip).
		Set("coredns_version", config.Cfg["coredns_version"]).
		SaveAs(coredns_yaml_sed_file).Close()
	//创建 deploy.sh
	new(common.FileAssist).
		CopyFrom(template_deploy_sh_file).
		Println("创建 "+deploy_sh_file_name+"...", ip).
		CopyTo(deploy_sh_file).Close()
	//创建 deploy.sh
	new(common.FileAssist).
		CopyFrom(template_dnsutils_yaml_file).
		Println("创建 "+dnsutils_yaml_file_name+"...", ip).
		CopyTo(dnsutils_yaml_file).Close()

}
