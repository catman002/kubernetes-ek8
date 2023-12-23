package assist

import (
	"ek8/common"
	"ek8/shell"
)

func Copy_cert_from_local_to_target(allca_dir string, sourceCertName string, CN string, targetCertName string, ip string) (state bool) {
	cmd := new(common.Command).
		Add("source config.cfg").
		//Add("rm -rf $release_dir/" + ip + "/cert/" + targetCertName).
		Add("mkdir -p $release_dir/" + ip + "/cert/" + targetCertName).
		Add("cp -Rf  " + allca_dir + "/" + sourceCertName + "/*" + " $release_dir/" + ip + "/cert/" + targetCertName).
		Add("echo \"证书 " + CN + " 拷贝完成\"").
		ToString()
	status, _ := shell.Exec(cmd, ip)
	return status
}
func Copy_certFile_from_local_to_targetFile(allca_dir string, sourceCertName string, target_dir string, targetCertName string, ip string) (state bool) {
	cmd := new(common.Command).
		Add("source config.cfg").
		//Add("rm -rf $release_dir/" + ip + "/cert/" + targetCertName).
		Add("mkdir -p $release_dir/" + ip + "/cert/" + target_dir).
		Add("cp -f  " + allca_dir + "/" + sourceCertName + " $release_dir/" + ip + "/cert/" + target_dir + "/" + targetCertName).
		Add("echo \"证书 " + targetCertName + " 拷贝完成\"").
		ToString()
	status, _ := shell.Exec(cmd, ip)
	return status
}
