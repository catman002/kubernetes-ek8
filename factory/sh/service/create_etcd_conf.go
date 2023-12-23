package service

import (
	"ek8/common"
	"ek8/config"
	"strconv"
)

func Create_etcd_conf(etcdServers map[string]string) {
	if common.ModFlags["etcdCONFIG"] != 1 {
		return
	}
	common.TPrintln("创建etcd配置...")

	var etcd_initial_cluster string
	i := 0
	for _, ip := range etcdServers {
		if etcd_initial_cluster == "" {
			etcd_initial_cluster = "etcd" + strconv.Itoa(i+1) + "=https://" + ip + ":" + config.Cfg["etcd_peer_port"] //2380
			//etcd_initial_cluster = nodeName + "=https://" + ip + ":" + config.Cfg["etcd_peer_port"] //2380
		} else {
			etcd_initial_cluster = etcd_initial_cluster + ",etcd" + strconv.Itoa(i+1) + "=https://" + ip + ":" + config.Cfg["etcd_peer_port"] //2380
			//etcd_initial_cluster = etcd_initial_cluster + "," + nodeName + "=https://" + ip + ":" + config.Cfg["etcd_peer_port"] //2380
		}
		i++
	}
	etcd_initial_cluster = "\"" + etcd_initial_cluster + "\""

	i = 0
	for _, ip := range etcdServers {
		target_dir := config.Cfg["release_dir"] + "/" + ip + "/service/etcd"

		etcd_conf_file_name := "etcd.conf"
		etcd_service_file_name := "etcd.service"
		etcd_conf_file := target_dir + "/" + etcd_conf_file_name
		etcd_service_file := target_dir + "/" + etcd_service_file_name

		template_etcd_conf_file := "factory/template/service/etcd/" + etcd_conf_file_name
		template_etcd_service_file := "factory/template/service/etcd/" + etcd_service_file_name

		etcd_name := "etcd" + strconv.Itoa(i+1)
		etcd_data_dir := config.Cfg["etcd_dir"] + "/data"

		etcd_listen_client_urls := "https://" + ip + ":" + config.Cfg["etcd_port"]    //2379
		etcd_advertise_client_urls := "https://" + ip + ":" + config.Cfg["etcd_port"] //2379

		etcd_listen_peer_urls := "https://" + ip + ":" + config.Cfg["etcd_peer_port"]            //2380
		etcd_initial_advertise_peer_urls := "https://" + ip + ":" + config.Cfg["etcd_peer_port"] //2380

		new(common.FileAssist).Load(template_etcd_conf_file).
			Println("创建 "+etcd_conf_file_name+"...", ip).
			Set("etcd_name", etcd_name).
			Set("etcd_data_dir", etcd_data_dir).
			Set("etcd_listen_client_urls", etcd_listen_client_urls).
			Set("etcd_advertise_client_urls", etcd_advertise_client_urls).
			Set("etcd_listen_peer_urls", etcd_listen_peer_urls).
			Set("etcd_initial_advertise_peer_urls", etcd_initial_advertise_peer_urls).
			Set("etcd_initial_cluster", etcd_initial_cluster).
			Set("ETCD_INITIAL_CLUSTER_STATE", config.ETCD_INITIAL_CLUSTER_STATE).
			SaveAs(etcd_conf_file).Close()
		new(common.FileAssist).Load(template_etcd_service_file).
			Println("创建 "+etcd_service_file_name+"...", ip).
			Set("etcd_cert_file", config.Cfg["etcd_cert_file"]).
			Set("etcd_key_file", config.Cfg["etcd_key_file"]).
			Set("etcd_trusted_ca_file", config.Cfg["etcd_trusted_ca_file"]).
			Set("etcd_peer_cert_file", config.Cfg["etcd_peer_cert_file"]).
			Set("etcd_peer_key_file", config.Cfg["etcd_peer_key_file"]).
			Set("etcd_peer_trusted_ca_file", config.Cfg["etcd_peer_trusted_ca_file"]).
			Set("etcd_dir", config.Cfg["etcd_dir"]).
			Set("etcd_exec_start", config.Cfg["etcd_exec_start"]).
			SaveAs(etcd_service_file).Close()
		i++
	}
}
