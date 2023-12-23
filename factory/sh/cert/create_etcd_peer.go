package cert

import (
	"ek8/common"
	"ek8/config"
	"ek8/factory/sh/cert/assist"
	"ek8/shell"
	"fmt"
	"sync"
)

func Create_etcd_peer(etcdServers map[string]string) {
	if !(common.ModFlags["etcdPeerCERT"] == 1 ||
		common.ModFlags["caCONFIG"] == 1 ||
		common.ModFlags["etcdCA"] == 1) {
		return
	}
	common.TPrintln("创建kube-etcd-peer证书 [peer]...")

	var once sync.Once
	for _, ip := range etcdServers {
		allca_dir := "target/" + ip + "/cert"
		etcd_dir_name := "etcd"
		etcd_dir := allca_dir + "/" + etcd_dir_name
		etcd_dir_tmp := temp_dir + "/" + etcd_dir_name

		//创建证书
		create_etcd_peer_cert(&once, ip, etcdServers, allca_dir, etcd_dir, etcd_dir_name, etcd_dir_tmp)
	}
	if config.Cfg["delete_cert_tem_dir"] == "1" {
		shell.Exec_simple("rm -rf  " + temp_dir)
	}

}

func create_etcd_peer_cert(once *sync.Once, ip string, etcdServers map[string]string,
	allca_dir string,
	etcd_dir string,
	etcd_dir_name string,
	etcd_dir_tmp string) {

	etcd_ca_CN := config.Cfg["etcd_ca_CN"]
	etcd_ca_file := etcd_dir + "/" + etcd_ca_CN + ".pem"
	etcd_key_file := etcd_dir + "/" + etcd_ca_CN + "-key.pem"
	ca_config_file := allca_dir + "/ca-config.json"

	etcd_peer_csr_file_name := "etcd-peer-csr.json"
	etcd_peer_csr_file := etcd_dir_tmp + "/" + etcd_peer_csr_file_name

	template_csr_json_file := "factory/template/cert/etcd/" + etcd_peer_csr_file_name
	CN := config.Cfg["etcd_peer_CN"]

	createCert := func() {
		new(common.FileAssist).
			Println("创建 "+etcd_peer_csr_file_name+".....", "local").
			Load(template_csr_json_file).
			Set("etcd_peer_CN", config.Cfg["etcd_peer_CN"]).
			SaveAs(etcd_peer_csr_file).Close()
		//创建etcd-server的证书请求、证书、私锁
		var etcd_hostname_list string
		for nodeName, ip := range etcdServers {
			etcd_hostname_list = etcd_hostname_list + ip + "," + nodeName + ","
		}
		etcd_hostname_list = etcd_hostname_list + "127.0.0.1,localhost"
		echo := "echo \"创建 %s.csr %s.pem %s-key.pem.......\""
		echo = fmt.Sprintf(echo, CN, CN, CN)
		cmd := new(common.Command).Add("source config.cfg").
			Add("cfssl gencert -ca=" + etcd_ca_file + " -ca-key=" + etcd_key_file + " -config=" + ca_config_file + " -hostname=" + etcd_hostname_list + " -profile=peer -loglevel=1 " + etcd_peer_csr_file + " | cfssljson -bare " + etcd_dir_tmp + "/" + CN).
			Add(echo).
			ToString()
		shell.Exec(cmd, "local")
	}
	once.Do(createCert)
	assist.Copy_cert_from_local_to_target(temp_dir, etcd_dir_name, CN, etcd_dir_name, ip)

	//删除临时文件
	shell.Exec_simple("rm -rf   " + etcd_dir + "/*.json " + etcd_dir + "/*.csr")
	shell.Exec_simple("rm -rf   " + etcd_dir_tmp + "/*.json " + etcd_dir_tmp + "/*.csr")

}
