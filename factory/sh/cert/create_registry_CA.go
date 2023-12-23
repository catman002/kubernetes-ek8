package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"strconv"
	"sync"
)

func Create_registry_CA(allservers map[string]string) {
	if !(common.ModFlags["caCONFIG"] == 1 ||
		common.ModFlags["registryCA"] == 1) {
		return
	}
	common.TPrintln("创建 registry-ca 证书 [ca]...")
	var once sync.Once
	for _, ip := range allservers {
		allca_dir := "$PWD/target/" + ip + "/cert"
		registry_dir_name := "registry"
		registry_dir := allca_dir + "/" + registry_dir_name
		registry_dir_temp := temp_dir + "/" + registry_dir_name

		registry_ca_csr_file_name := "registry-ca-csr.json"
		registry_ca_csr_file := registry_dir_temp + "/" + registry_ca_csr_file_name

		template_csr_json_file := "factory/template/cert/registry/" + registry_ca_csr_file_name
		CN := config.Cfg["registry_ca_CN"]
		ca_expire_time, _ := strconv.Atoi(config.CADAYS)
		ca_expire_time = ca_expire_time * 24

		createCert := func() {

			//创建registry-ca的证书请求的json配置文件
			new(common.FileAssist).
				Println("创建 "+registry_ca_csr_file_name+".....", "local").
				Load(template_csr_json_file).
				Set("registry_ca_CN", config.Cfg["registry_ca_CN"]).
				Set("ca_expire_time", strconv.Itoa(ca_expire_time)).
				SaveAs(registry_ca_csr_file).Close()
			//创建registry-ca的证书请求、证书、私锁
			echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
			echo = fmt.Sprintf(echo, CN, CN, CN)
			cmd := new(common.Command).
				Add("cfssl gencert -initca " + registry_ca_csr_file + " | cfssljson -bare " + registry_dir_temp + "/" + CN).
				Add(echo).
				ToString()
			shell.Exec(cmd, "local")
		}
		once.Do(createCert)
		assist.Copy_cert_from_local_to_target(temp_dir, registry_dir_name, CN, registry_dir_name, ip)

		//删除临时文件
		shell.Exec_simple("rm -rf   " + registry_dir + "/*.json " + registry_dir + "/*.csr")
		shell.Exec_simple("rm -rf   " + registry_dir_temp + "/*.json " + registry_dir_temp + "/*.csr")
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}
