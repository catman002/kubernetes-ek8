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

func Create_front_proxy_CA(apiservers map[string]string) {
	if !(common.ModFlags["frontProxyCA"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建ubernetes-front-proxy-ca证书 [ca]...")
	var once sync.Once
	for _, ip := range apiservers {
		allca_dir := "target/" + ip + "/cert"
		front_proxy_dir_name := "front-proxy"
		front_proxy_ca_dir := allca_dir + "/" + front_proxy_dir_name
		front_proxy_ca_dir_temp := temp_dir + "/" + front_proxy_dir_name

		template_csr_json_file := "factory/template/cert/front-proxy/front-proxy-ca-csr.json"

		front_proxy_ca_csr_file_name := "front-proxy-ca-csr.json"
		front_proxy_ca_csr_file := front_proxy_ca_dir_temp + "/" + front_proxy_ca_csr_file_name
		CN := config.Cfg["front_proxy_ca_CN"]
		ca_expire_time, _ := strconv.Atoi(config.CADAYS)
		ca_expire_time = ca_expire_time * 24

		createCert := func() {
			//创建etcd-ca的证书请求的json配置文件
			new(common.FileAssist).
				Println("创建 "+front_proxy_ca_csr_file_name+".....", "local").
				Load(template_csr_json_file).
				Set("front_proxy_ca_CN", config.Cfg["front_proxy_ca_CN"]).
				Set("ca_expire_time", strconv.Itoa(ca_expire_time)).
				SaveAs(front_proxy_ca_csr_file).Close()

			//创建etcd-ca的证书请求、证书、私锁
			echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
			echo = fmt.Sprintf(echo, CN, CN, CN)
			cmd := new(common.Command).
				Add("cfssl gencert -initca " + front_proxy_ca_csr_file + " | cfssljson -bare " + front_proxy_ca_dir_temp + "/" + CN).
				Add(echo).
				ToString()
			shell.Exec(cmd, "local")
		}
		once.Do(createCert)
		assist.Copy_cert_from_local_to_target(temp_dir, front_proxy_dir_name, CN, front_proxy_dir_name, ip)

		//删除临时文件
		shell.Exec_simple("rm -rf   " + front_proxy_ca_dir + "/*.json " + front_proxy_ca_dir + "/*.csr")
		shell.Exec_simple("rm -rf   " + front_proxy_ca_dir_temp + "/*.json " + front_proxy_ca_dir_temp + "/*.csr")

	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}
