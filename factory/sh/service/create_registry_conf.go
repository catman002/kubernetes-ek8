package service

import (
	"ek8/common"
	"ek8/config"
)

func Create_registry_conf() {
	if common.ModFlags["registryCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建registry配置...")
	target_dir := config.Cfg["release_dir"] + "/" + config.Cfg["registry_server_ip"] + "/service/registry"

	docker_registry_yaml_file_name := "docker-registry.yaml"
	docker_registry_yaml_file := target_dir + "/" + docker_registry_yaml_file_name

	template_docker_registry_yaml_file := "factory/template/service/registry/" + docker_registry_yaml_file_name

	//创建docker-registry.yaml
	new(common.FileAssist).
		Load(template_docker_registry_yaml_file).
		Println("创建 "+docker_registry_yaml_file_name+"...", config.Cfg["registry_server_ip"]).
		Set("registry_server_repo_url", config.Cfg["registry_server_repo_url"]).
		Set("registry_server_port", config.Cfg["registry_server_port"]).
		SaveAs(docker_registry_yaml_file).
		Close()

}
