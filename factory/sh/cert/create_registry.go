package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_registry(allServers map[string]string) {
	if !(common.ModFlags["registryCERT"] == 1 ||
		common.ModFlags["registryCA"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建registry证书 [server]...")
	var once sync.Once
	for _, ip := range allServers {
		allca_dir := "$PWD/target/" + ip + "/cert"
		registry_dir_name := "registry"
		registry_dir := allca_dir + "/" + registry_dir_name
		registry_dir_temp := temp_dir + "/" + registry_dir_name

		//创建证书
		create_registry_cert(&once, ip, allca_dir, registry_dir, registry_dir_name, registry_dir_temp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  target/.cert")
	}
}

func create_registry_cert(once *sync.Once,
	ip string,
	allca_dir string,
	registry_dir string,
	registry_dir_name string,
	registry_dir_temp string) {

	registry_ca_CN := config.Cfg["registry_ca_CN"]
	registry_ca_file := registry_dir + "/" + registry_ca_CN + ".pem"
	registry_key_file := registry_dir + "/" + registry_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	registry_csr_file_name := "registry-csr.json"
	registry_csr_file := registry_dir_temp + "/" + registry_csr_file_name

	template_csr_json_file := "factory/template/cert/registry/" + registry_csr_file_name
	CN := config.Cfg["registry_server_domain"] + ":" + config.Cfg["registry_server_port"]

	createCert := func() {
		//创建registry-server的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+registry_csr_file_name+".....", ip).
			Load(template_csr_json_file).
			Set("registry_CN", CN).
			SaveAs(registry_csr_file).Close()
		//创建registry-server的证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + registry_ca_file + " -ca-key=" + registry_key_file + " -config=" + ca_config_file + " -profile=server -loglevel=1 " + registry_csr_file + " | cfssljson -bare " + registry_dir_temp + "/registry").
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")

		//创建htpasswd
		cmd = new(common.Command).Add("source config.cfg").
			Add("htpasswd -Bbn " + config.Cfg["registry_username"] + " " + config.Cfg["registry_passwd"] + " > " + registry_dir_temp + "/htpasswd"). //密码认证文件
			Add("echo \"创建 htpasswd.....\"").
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, registry_dir_name, CN, registry_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + registry_dir + "/*.json " + registry_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + registry_dir_temp + "/*.json " + registry_dir_temp + "/*.csr")

}
