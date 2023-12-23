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

func Create_etcd_CA(etcdServers map[string]string) {
	if !(common.ModFlags["caCONFIG"] == 1 ||
		common.ModFlags["etcdCA"] == 1) {
		return
	}
	common.TPrintln("创建etcd-ca证书 [ca]...")
	var once sync.Once
	for _, ip := range etcdServers {
		allca_dir := "target/" + ip + "/cert"
		etcd_dir_name := "etcd"
		etcd_dir := allca_dir + "/" + etcd_dir_name
		etcd_dir_temp := temp_dir + "/" + etcd_dir_name

		etcd_ca_csr_file_name := "etcd-ca-csr.json"
		etcd_ca_csr_file := etcd_dir_temp + "/" + etcd_ca_csr_file_name

		template_csr_json_file := "factory/template/cert/etcd/" + etcd_ca_csr_file_name
		CN := config.Cfg["etcd_ca_CN"]
		ca_expire_time, _ := strconv.Atoi(config.CADAYS)
		ca_expire_time = ca_expire_time * 24

		createCert := func() {
			//创建etcd-ca的证书请求的json配置文件
			new(common.FileAssist).
				Println("创建 "+etcd_ca_csr_file_name+".....", "local").
				Load(template_csr_json_file).
				Set("etcd_ca_CN", config.Cfg["etcd_ca_CN"]).
				Set("ca_expire_time", strconv.Itoa(ca_expire_time)).
				SaveAs(etcd_ca_csr_file).Close()
			//创建etcd-ca的证书请求、证书、私锁
			echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
			echo = fmt.Sprintf(echo, CN, CN, CN)
			cmd := new(common.Command).
				Add("cfssl gencert -initca " + etcd_ca_csr_file + " | cfssljson -bare " + etcd_dir_temp + "/" + CN).
				Add(echo).
				ToString()
			shell.Exec(cmd, "local")
		}
		once.Do(createCert)
		assist.Copy_cert_from_local_to_target(temp_dir, etcd_dir_name, CN, etcd_dir_name, ip)

		//删除临时文件
		shell.Exec_simple("rm -rf   " + etcd_dir + "/*.json " + etcd_dir + "/*.csr")
		shell.Exec_simple("rm -rf   " + etcd_dir_temp + "/*.json " + etcd_dir_temp + "/*.csr")
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}
}
