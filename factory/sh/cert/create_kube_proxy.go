package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_kube_proxy(allServers map[string]string) {
	if !(common.ModFlags["kubeProxyCERT"] == 1 ||
		common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建kube-proxy证书 [client]...")

	var once sync.Once
	for _, ip := range allServers {
		allca_dir := "target/" + ip + "/cert"
		kube_proxy_dir_name := "kube-proxy"
		kube_proxy_dir := allca_dir + "/" + kube_proxy_dir_name
		kube_proxy_dir_temp := temp_dir + "/" + kube_proxy_dir_name

		//创建证书
		create_kube_proxy_cert(&once, ip, allca_dir, kube_proxy_dir, kube_proxy_dir_name, kube_proxy_dir_temp)

		//创建kubeconfig文件
		create_kube_proxy_sh_file(ip, kube_proxy_dir, kube_proxy_dir_temp)

	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}

}

func create_kube_proxy_cert(once *sync.Once,
	ip string,
	allca_dir string,
	kube_proxy_dir string,
	kube_proxy_dir_name string,
	kube_proxy_dir_temp string) {

	kubernetes_ca_CN := config.Cfg["ca_CN"]
	kubernetes_ca_file := allca_dir + "/ca/" + kubernetes_ca_CN + ".pem"
	kubernetes_key_file := allca_dir + "/ca/" + kubernetes_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	kube_proxy_csr_file_name := "kube-proxy-csr.json"
	kube_proxy_csr_file := kube_proxy_dir_temp + "/" + kube_proxy_csr_file_name

	template_csr_json_file := "factory/template/cert/kube-proxy/" + kube_proxy_csr_file_name
	CN := config.Cfg["kubeproxy_CN"]

	createCert := func() {
		//创建kube_apiserver-server的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+kube_proxy_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("kubeproxy_CN", config.Cfg["kubeproxy_CN"]).
			SaveAs(kube_proxy_csr_file).Close()

		//创建kube_apiserver-server的证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + kubernetes_ca_file + " -ca-key=" + kubernetes_key_file + " -config=" + ca_config_file + " -profile=client -loglevel=1 " + kube_proxy_csr_file + " | cfssljson -bare " + kube_proxy_dir_temp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, kube_proxy_dir_name, CN, kube_proxy_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + kube_proxy_dir + "/*.json " + kube_proxy_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + kube_proxy_dir_temp + "/*.json " + kube_proxy_dir_temp + "/*.csr")

}

func create_kube_proxy_sh_file(ip string,
	kube_proxy_dir string,
	kube_proxy_dir_temp string) {
	cluster_address := "https://" + config.Cfg["VIP"] + ":" + config.Cfg["VIP_PORT"]

	config_sh_file_name := "create_config.sh"
	config_sh_file := kube_proxy_dir + "/" + config_sh_file_name
	config_sh_file_temp := kube_proxy_dir_temp + "/" + config_sh_file_name

	template_config_sh_file := "factory/template/cert/kube-proxy/" + config_sh_file_name

	new(common.FileAssist).
		Println("创建 "+config_sh_file_name+".....", ip).
		Load(template_config_sh_file).
		Set("cluster_name", config.Cfg["cluster_name"]).
		Set("kubeproxy_certificate_authority", config.Cfg["kubeproxy_certificate_authority"]).
		Set("cluster_address", cluster_address).
		Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
		Set("kubeproxy_CN", config.Cfg["kubeproxy_CN"]).
		Set("kubeproxy_client_certificate", config.Cfg["kubeproxy_client_certificate"]).
		Set("kubeproxy_client_key", config.Cfg["kubeproxy_client_key"]).
		SaveAs(config_sh_file).
		SaveAs(config_sh_file_temp).
		Close()
}
