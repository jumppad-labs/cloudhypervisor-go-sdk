package sdk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kdomanski/iso9660"
)

func CreateCloudInitDisk(hostname string, password string, mac string, cidr string, gateway string, userdata string) (string, error) {
	source, err := os.MkdirTemp("", "cloudinit-*")
	if err != nil {
		return "", err
	}

	fmt.Println(source)

	err = os.WriteFile(filepath.Join(source, "meta-data"), []byte(fmt.Sprintf("instance-id: %s\nlocal-hostname: %s", hostname, hostname)), 0644)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(source, "user-data"), []byte(userdata), 0644)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(source, "network-config"), []byte(fmt.Sprintf(`version: 2
	ethernets:
	  ens4:
	    match:
	      macaddress: %s
	      addresses: [%s]
	      gateway4: %s
	`, mac, cidr, gateway)), 0644)
	if err != nil {
		return "", err
	}

	// check if machine is already running ... cant add disk then.
	// TODO: generate source files?
	destination, _ := filepath.Abs("/tmp/cloudinit.iso")
	err = createISO9660Disk(source, "cidata", destination)
	if err != nil {
		return "", err
	}

	return destination, nil
}

func createISO9660Disk(source string, label string, destination string) error {
	writer, err := iso9660.NewWriter()
	if err != nil {
		return err
	}
	defer writer.Cleanup()

	// err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if info.IsDir() {
	// 		return fmt.Errorf("directories are not supported")
	// 	}

	// 	relativePath, err := filepath.Rel(source, path)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	file, err := os.Open(path)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = writer.AddFile(file, relativePath)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }

	metadataSource := filepath.Join(source, "meta-data")
	mf, err := os.Open(metadataSource)
	if err != nil {
		return err
	}
	defer mf.Close()

	err = writer.AddFile(mf, "meta-data")
	if err != nil {
		return err
	}

	userdataSource := filepath.Join(source, "user-data")
	uf, err := os.Open(userdataSource)
	if err != nil {
		return err
	}
	defer uf.Close()

	err = writer.AddFile(uf, "user-data")
	if err != nil {
		return err
	}

	networkConfigSource := filepath.Join(source, "user-data")
	nf, err := os.Open(networkConfigSource)
	if err != nil {
		return err
	}
	defer nf.Close()

	err = writer.AddFile(nf, "network-config")
	if err != nil {
		return err
	}

	of, err := os.OpenFile(destination, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer of.Close()

	err = writer.WriteTo(of, label)
	if err != nil {
		return err
	}
	return nil
}
