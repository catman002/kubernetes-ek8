package service

import (
	"ek8/common"
	"ek8/config"
	"strings"
)

func Create_sh_for_init_env(allServers map[string]string) {
	if common.ModFlags["initEnvCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建环境配置...")

	for nodeName, ip := range allServers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/init-env"

		init_env_sh_file_name := "init_env.sh"
		init_hostname_sh_file_name := "init_hostname.sh"
		update_core_sh_file_name := "update_core.sh"
		init_env_sh_file := target_dir + "/" + init_env_sh_file_name
		init_hostname_sh_file := target_dir + "/" + init_hostname_sh_file_name
		update_core_sh_file := target_dir + "/" + config.Cfg["target_system"] + "/" + update_core_sh_file_name

		template_init_env_sh_file := "factory/template/service/init-env/" + init_env_sh_file_name
		template_init_hostname_sh_file := "factory/template/service/init-env/" + init_hostname_sh_file_name
		template_update_core_sh_file := "factory/template/service/init-env/" + config.Cfg["target_system"] + "/" + update_core_sh_file_name

		//创建init_env.sh
		hostname := nodeName
		new(common.FileAssist).
			CopyFrom(template_init_env_sh_file).
			Println("创建 "+init_env_sh_file_name+"...", ip).
			CopyTo(init_env_sh_file).
			Close()
		//创建init_hostname.sh
		new(common.FileAssist).
			Load(template_init_hostname_sh_file).
			Println("创建 "+init_hostname_sh_file_name+"...", ip).
			Set("hostname", hostname).
			Set("allServers", func() string {
				servers_str := ""
				for name, server := range allServers {
					servers_str = servers_str + name + ":" + server + " "
				}
				servers_str = strings.TrimRight(servers_str, " ")
				return servers_str
			}()).
			SaveAs(init_hostname_sh_file).
			Close()
		//创建update_core.sh
		new(common.FileAssist).
			CopyFrom(template_update_core_sh_file).
			Println("创建 "+update_core_sh_file_name+"...", ip).
			CopyTo(update_core_sh_file).
			Close()
	}
}
