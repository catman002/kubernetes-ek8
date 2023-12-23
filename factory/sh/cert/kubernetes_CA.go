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

func Create_kubernetes_CA(allservers map[string]string) {
	if !(common.ModFlags["ca"] == 1 ||
		common.ModFlags["caCONFIG"] == 1) {
		return
	}
	common.TPrintln("创建 kubernetes-ca证书 [ca]...")
	var once sync.Once
	for _, ip := range allservers {
		allca_dir := "$PWD/target/" + ip + "/cert"
		kubernetes_dir_name := "ca"
		kubernetes_dir := allca_dir + "/" + kubernetes_dir_name
		kubernetes_dir_tmp := temp_dir + "/" + kubernetes_dir_name

		kubernetes_ca_csr_file_name := "kube-ca-csr.json"
		kubernetes_ca_csr_file := kubernetes_dir_tmp + "/" + kubernetes_ca_csr_file_name

		template_csr_json_file := "factory/template/cert/ca/" + kubernetes_ca_csr_file_name
		CN := config.Cfg["ca_CN"]
		ca_expire_time, _ := strconv.Atoi(config.CADAYS)
		ca_expire_time = ca_expire_time * 24

		createCert := func() {
			//创建kubernetes-ca的证书请求的json配置文件
			new(common.FileAssist).
				Println("创建 "+kubernetes_ca_csr_file_name+".....", "local").
				Load(template_csr_json_file).
				Set("ca_CN", config.Cfg["ca_CN"]).
				Set("ca_expire_time", strconv.Itoa(ca_expire_time)).
				SaveAs(kubernetes_ca_csr_file).Close()

			//创建kubernetes-ca的证书请求、证书、私锁
			echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
			echo = fmt.Sprintf(echo, CN, CN, CN)
			cmd := new(common.Command).
				Add("cfssl gencert -initca " + kubernetes_ca_csr_file + " | cfssljson -bare " + kubernetes_dir_tmp + "/" + CN).
				Add(echo).
				ToString()
			shell.Exec(cmd, "local")
		}
		once.Do(createCert)
		assist.Copy_cert_from_local_to_target(temp_dir, kubernetes_dir_name, CN, kubernetes_dir_name, ip)

		//删除临时文件
		shell.Exec_simple("rm -rf   " + kubernetes_dir + "/*.json " + kubernetes_dir + "/*.csr")
		shell.Exec_simple("rm -rf   " + kubernetes_dir_tmp + "/*.json " + kubernetes_dir_tmp + "/*.csr")

	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}
