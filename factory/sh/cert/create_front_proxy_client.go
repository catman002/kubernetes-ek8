package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_front_proxy_client(apiservers map[string]string) {
	if !(common.ModFlags["frontProxyCA"] == 1 ||
		common.ModFlags["caCONFIG"] == 1 ||
		common.ModFlags["frontProxyClientCERT"] == 1) {
		return
	}
	common.TPrintln("创建front-proxy-client证书 [client]...")

	var once sync.Once
	for _, ip := range apiservers {
		allca_dir := "$PWD/target/" + ip + "/cert"
		front_proxy_dir_name := "front-proxy"
		front_proxy_dir := allca_dir + "/" + front_proxy_dir_name
		front_proxy_dir_temp := temp_dir + "/" + front_proxy_dir_name

		//创建证书
		create_front_proxy_client_cert(&once, ip, allca_dir, front_proxy_dir, front_proxy_dir_name, front_proxy_dir_temp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}

func create_front_proxy_client_cert(once *sync.Once, ip string,
	allca_dir string,
	front_proxy_dir string,
	front_proxy_dir_name string,
	front_proxy_dir_temp string) {

	front_proxy_ca_CN := config.Cfg["front_proxy_ca_CN"]
	front_proxy_ca_file := front_proxy_dir + "/" + front_proxy_ca_CN + ".pem"
	front_proxy_key_file := front_proxy_dir + "/" + front_proxy_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	front_proxy_client_csr_file_name := "front_proxy-client-csr.json"
	front_proxy_client_csr_file := front_proxy_dir_temp + "/" + front_proxy_client_csr_file_name
	template_csr_json_file := "factory/template/cert/front-proxy/" + front_proxy_client_csr_file_name
	CN := config.Cfg["frontproxyclient_CN"]

	createCert := func() {
		//创建etcd-client的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+front_proxy_client_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("frontproxyclient_CN", config.Cfg["frontproxyclient_CN"]).
			SaveAs(front_proxy_client_csr_file).Close()
		//创建etcd-client的证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + front_proxy_ca_file + " -ca-key=" + front_proxy_key_file + " -config=" + ca_config_file + " -profile=client -loglevel=1 " + front_proxy_client_csr_file + " | cfssljson -bare " + front_proxy_dir_temp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target("target/.cert", front_proxy_dir_name, CN, front_proxy_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + front_proxy_dir + "/*.json " + front_proxy_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + front_proxy_dir_temp + "/*.json " + front_proxy_dir_temp + "/*.csr")

}
