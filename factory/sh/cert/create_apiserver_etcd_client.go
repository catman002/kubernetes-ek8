package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_apiserver_etcd_client(allservers map[string]string) {
	if !(common.ModFlags["apiserverEtcdClientCERT"] == 1 ||
		common.ModFlags["caCONFIG"] == 1 ||
		common.ModFlags["etcdCA"] == 1) {
		return
	}
	common.TPrintln("创建kube-apiserver-etcd-client证书 [client]...")

	var once sync.Once
	for _, ip := range allservers {
		allca_dir := "target/" + ip + "/cert"
		apiserver_dir_name := "apiserver"
		apiserver_dir := allca_dir + "/" + apiserver_dir_name
		apiserver_dir_temp := temp_dir + "/" + apiserver_dir_name

		//创建证书
		create_apiserver_etcd_client_cert(&once, ip, allca_dir, apiserver_dir, apiserver_dir_name, apiserver_dir_temp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}

func create_apiserver_etcd_client_cert(once *sync.Once,
	ip string,
	allca_dir string,
	apiserver_dir string,
	apiserver_dir_name string,
	apiserver_dir_temp string) {

	etcd_ca_CN := config.Cfg["etcd_ca_CN"]
	etcd_ca_file := allca_dir + "/etcd/" + etcd_ca_CN + ".pem"
	etcd_ca_key_file := allca_dir + "/etcd/" + etcd_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	apiserver_etcd_client_csr_file_name := "apiserver-etcd-client-csr.json"
	apiserver_etcd_client_csr_file := apiserver_dir_temp + "/" + apiserver_etcd_client_csr_file_name

	template_csr_json_file := "factory/template/cert/apiserver/" + apiserver_etcd_client_csr_file_name
	CN := config.Cfg["apiserver_etcd_client_CN"]

	createCert := func() {
		//创建etcd-client的证书请求的json配置文件
		new(common.FileAssist).
			Println("创建 "+apiserver_etcd_client_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("apiserver_etcd_client_CN", config.Cfg["apiserver_etcd_client_CN"]).
			SaveAs(apiserver_etcd_client_csr_file).Close()
		//创建etcd-client的证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).
			Add("cfssl gencert -ca=" + etcd_ca_file + " -ca-key=" + etcd_ca_key_file + " -config=" + ca_config_file + " -profile=client -loglevel=1 " + apiserver_etcd_client_csr_file + " | cfssljson -bare " + apiserver_dir_temp + "/" + CN).
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
