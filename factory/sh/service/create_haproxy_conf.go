package service

import (
	"ek8/common"
	"ek8/config"
	"strconv"
)

func Create_haproxy_conf(haproxy map[string]string, apiServers map[string]string) {
	if common.ModFlags["haproxyCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建haproxy配置...")

	var back_servers string
	i := 0
	for _, ip := range apiServers {
		if back_servers == "" {
			back_servers = "server k8smaster" + strconv.Itoa(i+1) + " " + ip + ":" + config.Cfg["apiserver_secure_port"] + " maxconn 2000 weight 1 check inter 2000 rise 2 fall 2"
		} else {
			back_servers = back_servers + "\n  server k8smaster" + strconv.Itoa(i+1) + " " + ip + ":" + config.Cfg["apiserver_secure_port"] + " maxconn 2000 weight 1 check inter 2000 rise 2 fall 2"
		}
		i++
	}

	for _, ip := range haproxy {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/haproxy"

		haproxy_conf_file_name := "haproxy.conf"
		haproxy_service_file_name := "haproxy.service"
		haproxy_conf_file := target_dir + "/" + haproxy_conf_file_name
		haproxy_service_file := target_dir + "/" + haproxy_service_file_name

		template_haproxy_conf_file := "factory/template/service/haproxy/" + haproxy_conf_file_name
		template_haproxy_service_file := "factory/template/service/haproxy/" + haproxy_service_file_name

		//创建haproxy.service
		new(common.FileAssist).
			Load(template_haproxy_conf_file).
			Println("创建 "+haproxy_service_file_name+"...", ip).
			Set("VIP_PORT", config.Cfg["VIP_PORT"]).
			Set("apiserver_secure_port", config.Cfg["apiserver_secure_port"]).
			Set("back_servers", back_servers).
			SaveAs(haproxy_conf_file).
			Close()
		//创建haproxy.service
		new(common.FileAssist).
			Load(template_haproxy_service_file).
			Println("创建 "+haproxy_service_file_name+"...", ip).
			Set("haproxy_exec_start", config.Cfg["haproxy_exec_start"]).
			SaveAs(haproxy_service_file).
			Close()
	}

}
