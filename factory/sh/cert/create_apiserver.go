package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"strings"
	"sync"
)

func Create_apiserver(apiservers map[string]string) {
	if !(common.ModFlags["apiserverCERT"] == 1 ||
		common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建kube-apiserver证书 [peer]...")

	var once sync.Once
	for _, ip := range apiservers {
		allca_dir := "target/" + ip + "/cert"
		kube_apiserver_dir_name := "apiserver"
		kube_apiserver_dir := allca_dir + "/" + kube_apiserver_dir_name
		kube_apiserver_dir_tmp := temp_dir + "/" + kube_apiserver_dir_name

		//创建证书
		create_apiserver_cert(&once, ip, apiservers, allca_dir, kube_apiserver_dir, kube_apiserver_dir_name, kube_apiserver_dir_tmp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}

}

func create_apiserver_cert(once *sync.Once,
	ip string,
	allservers map[string]string,
	allca_dir string,
	kube_apiserver_dir string,
	kube_apiserver_dir_name string,
	kube_apiserver_dir_tmp string) {

	kubernetes_ca_CN := config.Cfg["ca_CN"]
	kubernetes_ca_file := allca_dir + "/ca/" + kubernetes_ca_CN + ".pem"
	kubernetes_key_file := allca_dir + "/ca/" + kubernetes_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	kube_apiserver_csr_file_name := "kube-apiserver-csr.json"
	template_csr_json_file := "factory/template/cert/apiserver/" + kube_apiserver_csr_file_name
	kube_apiserver_csr_file := kube_apiserver_dir_tmp + "/" + kube_apiserver_csr_file_name
	CN := config.Cfg["apiserver_CN"]

	createCert := func() {
		//创建kube_apiserver-server的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+kube_apiserver_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("apiserver_CN", config.Cfg["apiserver_CN"]).
			SaveAs(kube_apiserver_csr_file).Close()
		//创建kube_apiserver-server的证书请求、证书、私锁
		service_cluster_ip_range := config.Cfg["service_cluster_ip_range"]
		var kube_apiserver_hostname_list string
		for nodeName, _ip := range allservers {
			kube_apiserver_hostname_list = kube_apiserver_hostname_list + _ip + "," + nodeName + ","
		}
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "localhost,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "127.0.0.1,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "kubernetes.default.svc.cluster.loccal,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "kubernetes.default.svc.cluster,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "kubernetes.default.svc,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "kubernetes.default,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "kubernetes,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + config.Cfg["VIP"] + ","
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + service_cluster_ip_range[0:strings.LastIndex(service_cluster_ip_range, ".")] + ".1,"
		kube_apiserver_hostname_list = kube_apiserver_hostname_list + "istiod.istio-system.svc"
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + kubernetes_ca_file + " -ca-key=" + kubernetes_key_file + " -config=" + ca_config_file + " -hostname=" + kube_apiserver_hostname_list + " -profile=peer -loglevel=1 " + kube_apiserver_csr_file + " | cfssljson -bare " + kube_apiserver_dir_tmp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, kube_apiserver_dir_name, CN, kube_apiserver_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + kube_apiserver_dir + "/*.json " + kube_apiserver_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + kube_apiserver_dir_tmp + "/*.json " + kube_apiserver_dir_tmp + "/*.csr")

}
