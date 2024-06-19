package sdk

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kdomanski/iso9660"
)

//go:embed configs/meta-data.tmpl
var metadata string

//go:embed configs/user-data.tmpl
var userdata string

//go:embed configs/network-config.tmpl
var networkConfig string

func CreateCloudInitDisk(hostname string, mac string, cidr string, gateway string, username string, password string) (string, error) {
	source, err := os.MkdirTemp("", "cloudinit-*")
	if err != nil {
		return "", err
	}

	fmt.Println(source)

	err = generateMetadata(source, hostname)
	if err != nil {
		return "", err
	}

	err = generateUserdata(source, username, password)
	if err != nil {
		return "", err
	}

	err = generateNetworkConfig(source, mac, cidr, gateway)
	if err != nil {
		return "", err
	}

	// check if machine is already running ... cant add disk then.
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

	networkConfigSource := filepath.Join(source, "network-config")
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

func generateMetadata(destination string, hostname string) error {
	tmpl, err := template.New("meta-data").Parse(metadata)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(destination, "meta-data"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{
		"hostname": hostname,
	})
	if err != nil {
		return err
	}

	return nil
}

func generateUserdata(destination string, username string, password string) error {
	tmpl, err := template.New("user-data").Parse(userdata)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(destination, "user-data"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}

	return nil
}

func generateNetworkConfig(destination string, mac string, cidr string, gateway string) error {
	tmpl, err := template.New("network-config").Parse(networkConfig)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(destination, "network-config"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{
		"interface": "eth0",
		"mac":       mac,
		"cidr":      cidr,
		"gateway":   gateway,
	})
	if err != nil {
		return err
	}

	return nil
}
