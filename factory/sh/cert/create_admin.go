package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_admin(apiServers map[string]string) {
	if !(common.ModFlags["adminCERT"] == 1 ||
		common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建admin证书 [client]...")
	var once sync.Once

	for _, ip := range apiServers {
		allca_dir := "target/" + ip + "/cert"
		kube_admin_dir_name := "admin"
		kube_admin_dir := allca_dir + "/" + kube_admin_dir_name
		kube_admin_dir_temp := temp_dir + "/" + kube_admin_dir_name

		//创建证书
		create_admin_cert(&once, ip, allca_dir, kube_admin_dir, kube_admin_dir_name, kube_admin_dir_temp)

		//创建 create_clusterrolebinding.sh文件
		create_admin_cluster_role_binding_sh_file(ip, kube_admin_dir, kube_admin_dir_temp)

		//创建kubeconfig文件
		create_admin_kubeconfig(ip, kube_admin_dir, kube_admin_dir_temp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}

}

func create_admin_cert(once *sync.Once,
	ip string,
	allca_dir string,
	kube_admin_dir string,
	kube_admin_dir_name string,
	kube_admin_dir_temp string) {

	kubernetes_ca_CN := config.Cfg["ca_CN"]
	kubernetes_ca_file := allca_dir + "/ca/" + kubernetes_ca_CN + ".pem"
	kubernetes_key_file := allca_dir + "/ca/" + kubernetes_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	kube_admin_csr_json_file_name := "admin-csr.json"
	kube_admin_csr_json_file := kube_admin_dir_temp + "/" + kube_admin_csr_json_file_name

	template_csr_json_file := "factory/template/cert/admin/" + kube_admin_csr_json_file_name
	CN := config.Cfg["admin_CN"]

	createCert := func() {
		//创建son配置文件
		new(common.FileAssist).
			Println("创建 "+kube_admin_csr_json_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("admin_CN", config.Cfg["admin_CN"]).
			SaveAs(kube_admin_csr_json_file).Close()
		//创建kube-admin的证书请求、证书、私锁
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + kubernetes_ca_file + " -ca-key=" + kubernetes_key_file + " -config=" + ca_config_file + " -profile=client -loglevel=1 " + kube_admin_csr_json_file + " | cfssljson -bare " + kube_admin_dir_temp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, kube_admin_dir_name, CN, kube_admin_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + kube_admin_dir + "/*.json " + kube_admin_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + kube_admin_dir_temp + "/*.json " + kube_admin_dir_temp + "/*.csr")

}

func create_admin_kubeconfig(ip string,
	kube_admin_dir string,
	kube_admin_dir_temp string) {
	cluster_address := "https://" + config.Cfg["VIP"] + ":" + config.Cfg["VIP_PORT"]

	config_sh_file_name := "create_config.sh"
	config_sh_file := kube_admin_dir + "/" + config_sh_file_name
	config_sh_file_temp := kube_admin_dir_temp + "/" + config_sh_file_name

	template_config_sh_file := "factory/template/cert/admin/" + config_sh_file_name
	//创建 kubeconfig file
	new(common.FileAssist).
		Println("创建 "+config_sh_file_name+".....", ip).
		Load(template_config_sh_file).
		Set("cluster_name", config.Cfg["cluster_name"]).
		Set("kubectl_certificate_authority", config.Cfg["kubectl_certificate_authority"]).
		Set("cluster_address", cluster_address).
		Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
		Set("admin_CN", config.Cfg["admin_CN"]).
		Set("kubectl_client_certificate", config.Cfg["kubectl_client_certificate"]).
		Set("kubectl_client_key", config.Cfg["kubectl_client_key"]).
		SaveAs(config_sh_file).
		SaveAs(config_sh_file_temp).
		Close()

}

func create_admin_cluster_role_binding_sh_file(ip string,
	kube_admin_dir string,
	kube_admin_dir_temp string) {
	create_clusterrolebinding_sh_file_name := "create_clusterrolebinding.sh"
	create_clusterrolebinding_sh_file := kube_admin_dir + "/" + create_clusterrolebinding_sh_file_name
	create_clusterrolebinding_sh_file_temp := kube_admin_dir_temp + "/" + create_clusterrolebinding_sh_file_name

	template_clusterrolebinding_sh_file := "factory/template/cert/admin/" + create_clusterrolebinding_sh_file_name

	//创建clusterrolebinding.sh
	new(common.FileAssist).
		Println("创建 "+create_clusterrolebinding_sh_file_name+".....", ip).
		Load(template_clusterrolebinding_sh_file).
		Set("admin_CN", config.Cfg["admin_CN"]).
		Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
		SaveAs(create_clusterrolebinding_sh_file).
		SaveAs(create_clusterrolebinding_sh_file_temp).Close()
}
