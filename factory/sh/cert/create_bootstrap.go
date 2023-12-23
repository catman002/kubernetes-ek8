package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/shell"
	"os"
	"strings"
)

func Create_bootstrap(allservers map[string]string) {
	if !(common.ModFlags["bootstrapCERT"] == 1 ||
		common.ModFlags["tokenCSV"] == 1 ||
		common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建bootstrap [client]...")

	for _, ip := range allservers {
		allca_dir := "target/" + ip + "/cert"
		kube_bootstrap_dir_name := "bootstrap"
		kube_bootstrap_dir := allca_dir + "/" + kube_bootstrap_dir_name
		kube_bootstrap_dir_tmp := temp_dir + "/" + kube_bootstrap_dir_name

		//创建目录
		cmd := new(common.Command).
			Add("source config.cfg").
			Add("mkdir -p " + kube_bootstrap_dir).
			ToString()
		shell.Exec(cmd, "local")

		//创建kubeconfig文件
		create_bootstrap_kubeconfig_sh_file(ip, kube_bootstrap_dir, kube_bootstrap_dir+"/token.csv", kube_bootstrap_dir_tmp)

		//创建 tls_instructs_csr.yaml文件
		create_tls_instructs_csr_yaml_file(ip, kube_bootstrap_dir, kube_bootstrap_dir_tmp)

		//创建 create_clusterrolebinding.sh文件
		create_bootstrap_cluster_role_binding_sh_file(ip, kube_bootstrap_dir, kube_bootstrap_dir_tmp)

	}

}
func create_bootstrap_kubeconfig_sh_file(ip string,
	kube_bootstrap_dir string,
	token_file string,
	kube_bootstrap_dir_tmp string) {

	cluster_address := "https://" + config.Cfg["VIP"] + ":" + config.Cfg["VIP_PORT"]

	config_sh_file_name := "create_config.sh"
	config_sh_file := kube_bootstrap_dir + "/" + config_sh_file_name
	config_sh_file_tmp := kube_bootstrap_dir_tmp + "/" + config_sh_file_name

	template_config_sh_file := "factory/template/cert/bootstrap/" + config_sh_file_name

	//读取token
	token_line, _ := os.ReadFile(token_file)
	bootstrap_token := string(token_line)
	bootstrap_token = strings.Split(bootstrap_token, ",")[0]

	new(common.FileAssist).
		Println("创建 "+config_sh_file_name+".....", ip).
		Load(template_config_sh_file).
		Set("cluster_name", config.Cfg["cluster_name"]).
		Set("ca_file", config.Cfg["ca_file"]).
		Set("cluster_address", cluster_address).
		Set("kubernetes_etc_dir", config.Cfg["kubernetes_etc_dir"]).
		Set("bootstrap_user", config.Cfg["bootstrap_user"]).
		Set("bootstrap_token", bootstrap_token).
		SaveAs(config_sh_file).
		SaveAs(config_sh_file_tmp).
		Close()
}

func create_tls_instructs_csr_yaml_file(ip string,
	kube_bootstrap_dir string,
	kube_bootstrap_dir_tmp string) {
	tls_instructs_csr_yaml_file_name := "tls_instructs_csr.yaml"
	tls_instructs_csr_yaml_file := kube_bootstrap_dir + "/" + tls_instructs_csr_yaml_file_name
	tls_instructs_csr_yaml_file_tmp := kube_bootstrap_dir_tmp + "/" + tls_instructs_csr_yaml_file_name

	template_tls_instructs_csr_yaml_file := "factory/template/cert/bootstrap/" + tls_instructs_csr_yaml_file_name

	new(common.FileAssist).
		Println("创建 "+tls_instructs_csr_yaml_file_name+".....", ip).
		Load(template_tls_instructs_csr_yaml_file).
		SaveAs(tls_instructs_csr_yaml_file).
		SaveAs(tls_instructs_csr_yaml_file_tmp).
		Close()
}

func create_bootstrap_cluster_role_binding_sh_file(ip string,
	kube_bootstrap_dir string,
	kube_bootstrap_dir_tmp string) {
	create_clusterrolebinding_sh_file_name := "create_clusterrolebinding.sh"
	create_clusterrolebinding_sh_file := kube_bootstrap_dir + "/" + create_clusterrolebinding_sh_file_name
	create_clusterrolebinding_sh_file_tmp := kube_bootstrap_dir_tmp + "/" + create_clusterrolebinding_sh_file_name

	template_create_clusterrolebinding_sh_file := "factory/template/cert/bootstrap/" + create_clusterrolebinding_sh_file_name

	new(common.FileAssist).
		Println("创建 "+create_clusterrolebinding_sh_file_name+".....", ip).
		Load(template_create_clusterrolebinding_sh_file).
		Set("bootstrap_group", config.Cfg["bootstrap_group"]).
		Set("bootstrap_user", config.Cfg["bootstrap_user"]).
		SaveAs(create_clusterrolebinding_sh_file).
		SaveAs(create_clusterrolebinding_sh_file_tmp).
		Close()
}
