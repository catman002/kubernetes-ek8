package service

import (
	"ek8/common"
	"ek8/config"
)

func Create_calico_conf(allServices map[string]string) {
	if common.ModFlags["calicoCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建calico配置...")

	for _, ip := range common.MapToArray(allServices) {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/calico"

		calico_yaml_file_name := "calico.yaml"
		delete_net_sh_file_name := "delete_net.sh"
		calico_yaml_file := target_dir + "/" + calico_yaml_file_name
		delete_net_sh_file := target_dir + "/" + delete_net_sh_file_name

		template_calico_yaml_file := "factory/template/service/calico/" + config.Cfg["calico_version"] + "/" + calico_yaml_file_name
		template_delete_net_sh_file := "factory/template/service/calico/" + config.Cfg["calico_version"] + "/" + delete_net_sh_file_name

		new(common.FileAssist).
			Load(template_calico_yaml_file).
			Println("创建 "+calico_yaml_file_name+"...", ip).
			Set("interface", config.Cfg["interface"]).
			Set("cluster_cidr", config.Cfg["cluster_cidr"]).
			SaveAs(calico_yaml_file).Close()
		//创建 delete_net.sh
		new(common.FileAssist).
			CopyFrom(template_delete_net_sh_file).
			Println("创建 "+delete_net_sh_file_name+"...", ip).
			CopyTo(delete_net_sh_file).Close()

	}

}
