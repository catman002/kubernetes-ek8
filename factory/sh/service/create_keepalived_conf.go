package service

import (
	"ek8/common"
	"ek8/config"
	"strconv"
)

func Create_keepalived_conf(haproxy map[string]string) {
	if common.ModFlags["keepalivedCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建keepalived配置...")

	i := 0
	priority, _ := strconv.Atoi(config.Cfg["priority"])
	for _, ip := range haproxy {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/keepalived"

		keepalived_conf_file_name := "keepalived.conf"
		keepalived_service_file_name := "keepalived.service"
		keepalived_conf_file := target_dir + "/" + keepalived_conf_file_name
		keepalived_service_file := target_dir + "/" + keepalived_service_file_name

		check_server_sh_file_name := "check_server.sh"
		start_sh_file_name := "start.sh"
		locale_md_file_name := "locale.md"
		check_server_sh_file := target_dir + "/script/" + check_server_sh_file_name
		start_sh_file := target_dir + "/script/" + start_sh_file_name
		locale_md_file := target_dir + "/script/" + locale_md_file_name

		template_keepalived_conf_file := "factory/template/service/keepalived/" + keepalived_conf_file_name
		template_keepalived_service_file := "factory/template/service/keepalived/" + keepalived_service_file_name
		template_locale_md_file := "factory/template/service/keepalived/script/" + locale_md_file_name
		template_start_sh_file := "factory/template/service/keepalived/script/" + start_sh_file_name
		template_check_server_sh_file := "factory/template/service/keepalived/script/" + check_server_sh_file_name

		state := "MASTER"
		if i == 0 {
			state = "MASTER"
		} else {
			state = "BACKUP"
		}
		priority = priority - 1

		//创建keepalived.service
		new(common.FileAssist).
			Load(template_keepalived_service_file).
			Println("创建 "+keepalived_service_file_name+"...", ip).
			Set("keepalived_exec_start", config.Cfg["keepalived_exec_start"]).
			SaveAs(keepalived_service_file).
			Close()
		//创建keepalived.conf
		new(common.FileAssist).
			Load(template_keepalived_conf_file).
			Println("创建 "+keepalived_conf_file_name+"...", ip).
			Set("interface", config.Cfg["interface"]).
			Set("mcast_group", config.Cfg["mcast_group"]).
			Set("VIP", config.Cfg["VIP"]).
			Set("VIP_PORT", config.Cfg["VIP_PORT"]).
			Set("apiserver_secure_port", config.Cfg["apiserver_secure_port"]).
			Set("state", state).
			Set("netmask_bit", config.Cfg["netmask_bit"]).
			Set("weight", config.Cfg["weight"]).
			Set("vrid", config.Cfg["vrid"]).
			Set("priority", strconv.Itoa(priority)).
			Set("RID", "LVS"+strconv.Itoa(i)).
			SaveAs(keepalived_conf_file).
			Close()
		//创建check_server.sh
		new(common.FileAssist).
			CopyFrom(template_check_server_sh_file).
			Println("创建 "+check_server_sh_file_name+"...", ip).
			CopyTo(check_server_sh_file).
			Close()
		//创建start.sh
		new(common.FileAssist).
			CopyFrom(template_start_sh_file).
			Println("创建 "+start_sh_file_name+"...", ip).
			CopyTo(start_sh_file).
			Close()
		//创建locale.md
		new(common.FileAssist).
			CopyFrom(template_locale_md_file).
			Println("创建 "+locale_md_file_name+"...", ip).
			CopyTo(locale_md_file).
			Close()
		i++
	}
}
