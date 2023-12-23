package service

import (
	"ek8/common"
	"ek8/config"
)

func Create_docker_conf(allServers map[string]string) {
	if common.ModFlags["dockerCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建docker配置...")
	for _, ip := range allServers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/docker"

		docker_service_file_name := "docker.service"
		daemon_json_file_name := "daemon.json"
		delete_docker_var_file_name := "delete_docker_var.sh"

		docker_service_file := target_dir + "/" + docker_service_file_name
		daemon_json_file := target_dir + "/" + daemon_json_file_name
		delete_docker_var_file := target_dir + "/" + delete_docker_var_file_name

		template_daemon_json_file := "factory/template/service/docker/" + daemon_json_file_name
		template_docker_service_file := "factory/template/service/docker/" + docker_service_file_name
		template_delete_docker_var_file := "factory/template/service/docker/" + delete_docker_var_file_name

		//创建docker.service
		new(common.FileAssist).
			Load(template_docker_service_file).
			Println("创建 "+docker_service_file_name+"...", ip).
			Set("docker_data_dir", config.Cfg["docker_data_dir"]).
			Set("docker_work_dir", config.Cfg["docker_work_dir"]).
			Set("docker_exec_start", config.Cfg["docker_exec_start"]).
			SaveAs(docker_service_file).
			Close()
		//创建delete_docker_var.sh
		new(common.FileAssist).
			Println("创建 "+delete_docker_var_file_name+"...", ip).
			CopyFrom(template_delete_docker_var_file).
			CopyTo(delete_docker_var_file).
			Close()
		//创建daemon.json
		var registry_server = ""
		if config.Cfg["registry_server_domain"] != "" && config.Cfg["registry_server_port"] != "" {
			registry_server = config.Cfg["registry_server_domain"] + ":" + config.Cfg["registry_server_port"]
		}
		new(common.FileAssist).
			Println("创建 "+daemon_json_file_name+"...", ip).
			Load(template_daemon_json_file).
			Set("registry_server", registry_server).
			SaveAs(daemon_json_file).
			Close()
	}
}
