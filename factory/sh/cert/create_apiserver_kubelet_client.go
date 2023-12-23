package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_apiserver_kubelet_client(allservers map[string]string) {
	if !(common.ModFlags["apiserverKubeletClientCERT"] == 1 ||
		common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建kube-apiserver-kubelet-client证书 [client]...")

	var once sync.Once
	for _, ip := range allservers {
		allca_dir := "target/" + ip + "/cert"
		apiserver_dir_name := "apiserver"
		apiserver_dir := allca_dir + "/" + apiserver_dir_name
		apiserver_dir_temp := temp_dir + "/" + apiserver_dir_name

		//创建证书
		create_apiserver_kublet_client_cert(&once, ip, allca_dir, apiserver_dir, apiserver_dir_name, apiserver_dir_temp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}

func create_apiserver_kublet_client_cert(once *sync.Once,
	ip string,
	allca_dir string,
	apiserver_dir string,
	apiserver_dir_name string,
	apiserver_dir_temp string) {

	kubernetes_ca_CN := "kubernetes-ca"
	kubernetes_ca_file := allca_dir + "/ca/" + kubernetes_ca_CN + ".pem"
	kubernetes_key_file := allca_dir + "/ca/" + kubernetes_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	apiserver_kubelet_client_csr_file_name := "apiserver-kubelet-client-csr.json"
	apiserver_kubelet_client_csr_file := apiserver_dir_temp + "/" + apiserver_kubelet_client_csr_file_name

	template_csr_json_file := "factory/template/cert/apiserver/" + apiserver_kubelet_client_csr_file_name
	CN := config.Cfg["apiserver_kubelet_client_CN"]

	createCert := func() {
		//创建etcd-client的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+apiserver_kubelet_client_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("apiserver_kubelet_client_CN", config.Cfg["apiserver_kubelet_client_CN"]).
			SaveAs(apiserver_kubelet_client_csr_file).Close()
		//创建apiserver-kublet-client的证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + kubernetes_ca_file + " -ca-key=" + kubernetes_key_file + " -config=" + ca_config_file + " -profile=client -loglevel=1 " + apiserver_kubelet_client_csr_file + " | cfssljson -bare " + apiserver_dir_temp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, apiserver_dir_name, CN, apiserver_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + apiserver_dir + "/*.json " + apiserver_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + apiserver_dir_temp + "/*.json " + apiserver_dir_temp + "/*.csr")

}
