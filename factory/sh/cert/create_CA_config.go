package cert

import (
	"ek8/common"
	"ek8/config"
	"strconv"
)

var temp_dir = "target/.cert"

func Create_CA_config(allservers map[string]string) {
	if common.ModFlags["caCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建 ca-config.json [ca]...")

	for _, ip := range allservers {
		allca_dir := "target/" + ip + "/cert"

		ca_config_file_name := "ca-config.json"
		ca_config_file := allca_dir + "/" + ca_config_file_name

		template_ca_config_json_file := "factory/template/cert/" + ca_config_file_name

		ca_expire_time, _ := strconv.Atoi(config.CADAYS)
		ca_expire_time = ca_expire_time * 24

		//创建etcd-ca的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+ca_config_file_name+"...", ip).
			Load(template_ca_config_json_file).
			Set("ca_expire_time", strconv.Itoa(ca_expire_time)).
			SaveAs(ca_config_file).Close()
	}
}
