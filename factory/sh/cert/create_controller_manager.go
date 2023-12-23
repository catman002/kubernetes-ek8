package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_controller_manager(apiServer map[string]string) {
	if !(common.ModFlags["controllerManagerCERT"] == 1 ||
		common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建controller-manager证书 [client]...")

	var once sync.Once
	for _, ip := range apiServer {
		allca_dir := "target/" + ip + "/cert"
		kube_controller_manager_dir_name := "controller-manager"
		kube_controller_manager_dir := allca_dir + "/" + kube_controller_manager_dir_name
		kube_controller_manager_dir_temp := temp_dir + "/" + kube_controller_manager_dir_name
		//创建证书
		create_controller_manager_cert(&once, ip, allca_dir, kube_controller_manager_dir, kube_controller_manager_dir_name, kube_controller_manager_dir_temp)

		//创建kubeconfig文件
		create_controller_manager_sh_file(ip, kube_controller_manager_dir, kube_controller_manager_dir_temp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}

func create_controller_manager_cert(once *sync.Once,
	ip string,
	allca_dir string,
	kube_controller_manager_dir string,
	kube_controller_manager_dir_name string,
	kube_controller_manager_dir_temp string) {

	kubernetes_ca_CN := config.Cfg["ca_CN"]
	kubernetes_ca_file := allca_dir + "/ca/" + kubernetes_ca_CN + ".pem"
	kubernetes_key_file := allca_dir + "/ca/" + kubernetes_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	kube_controller_manager_csr_file_name := "kube-controller-manager-csr.json"
	kube_controller_manager_csr_file := kube_controller_manager_dir_temp + "/" + kube_controller_manager_csr_file_name

	template_csr_json_file := "factory/template/cert/controller-manager/" + kube_controller_manager_csr_file_name
	CN := config.Cfg["controller_CN"]

	createCert := func() {
		//创建kube_apiserver-server的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+kube_controller_manager_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("controller_CN", config.Cfg["controller_CN"]).
			SaveAs(kube_controller_manager_csr_file).Close()

		//创建证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + kubernetes_ca_file + " -ca-key=" + kubernetes_key_file + " -config=" + ca_config_file + " -profile=client -loglevel=1 " + kube_controller_manager_csr_file + " | cfssljson -bare " + kube_controller_manager_dir_temp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, kube_controller_manager_dir_name, CN, kube_controller_manager_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + kube_controller_manager_dir + "/*.json " + kube_controller_manager_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + kube_controller_manager_dir_temp + "/*.json " + kube_controller_manager_dir_temp + "/*.csr")

}

func create_controller_manager_sh_file(
	ip string,
	kube_controller_manager_dir string,
	kube_controller_manager_dir_temp string) {
	cluster_address := "https://" + config.Cfg["VIP"] + ":" + config.Cfg["VIP_PORT"]

	config_sh_file_name := "create_config.sh"
	config_sh_file := kube_controller_manager_dir + "/" + config_sh_file_name

	template_config_sh_file := "factory/template/cert/controller-manager/" + config_sh_file_name
	template_config_sh_file_temp := kube_controller_manager_dir_temp + "/" + config_sh_file_name

	//创建 kubeconfig file
	new(common.FileAssist).
		Println("创建 "+config_sh_file_name+".....", ip).
		Load(template_config_sh_file).
		Set("cluster_name", config.Cfg["cluster_name"]).
		Set("controller_certificate_authority", config.Cfg["controller_certificate_authority"]).
		Set("cluster_address", cluster_address).
		Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
		Set("controller_CN", config.Cfg["controller_CN"]).
		Set("controller_client_certificate", config.Cfg["controller_client_certificate"]).
		Set("controller_client_key", config.Cfg["controller_client_key"]).
		SaveAs(config_sh_file).
		SaveAs(template_config_sh_file_temp).
		Close()
}
