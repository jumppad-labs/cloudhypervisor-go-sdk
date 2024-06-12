package sdk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

func folderSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func createISO9660Disk(source string, label string, destination string) error {
	folderSize, err := folderSize(source)
	if err != nil {
		return err
	}

	var LogicalBlocksize diskfs.SectorSize = 2048
	mydisk, err := diskfs.Create(destination, folderSize, diskfs.Raw, LogicalBlocksize)
	if err != nil {
		return err
	}

	fspec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: label,
	}
	fs, err := mydisk.CreateFilesystem(fspec)
	if err != nil {
		return err
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = fs.Mkdir(relPath)
			if err != nil {
				return err
			}
			return nil
		}

		if !info.IsDir() {
			rw, err := fs.OpenFile(relPath, os.O_CREATE|os.O_RDWR)
			if err != nil {
				return err
			}

			in, errorOpeningFile := os.Open(path)
			if errorOpeningFile != nil {
				return errorOpeningFile
			}
			defer in.Close()

			_, err = io.Copy(rw, in)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	iso, ok := fs.(*iso9660.FileSystem)
	if !ok {
		if err != nil {
			return fmt.Errorf("not an iso9660 filesystem")
		}
	}
	err = iso.Finalize(iso9660.FinalizeOptions{})
	if err != nil {
		return err
	}

	return nil
}

// func createCloudInitDisk(source string) (string, error) {
// 	writer, err := iso9660.NewWriter()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer writer.Cleanup()

// 	output, err := os.MkdirTemp("", "cloudinit-*")
// 	if err != nil {
// 		return "", err
// 	}

// 	metadataSource := path.Join(source, "meta-data")
// 	metadata, err := os.ReadFile(metadataSource)
// 	if err != nil {
// 		return "", err
// 	}

// 	metadataDestination := filepath.Join(output, "meta-data")
// 	err = os.WriteFile(metadataDestination, []byte(metadata), 0644)
// 	if err != nil {
// 		return "", err
// 	}

// 	mf, err := os.Open(metadataDestination)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer mf.Close()

// 	err = writer.AddFile(mf, "meta-data")
// 	if err != nil {
// 		return "", err
// 	}

// 	userdataSource := path.Join(source, "user-data")
// 	userdata, err := os.ReadFile(userdataSource)
// 	if err != nil {
// 		return "", err
// 	}

// 	userdataDestination := path.Join(output, "user-data")
// 	err = os.WriteFile(userdataDestination, []byte(userdata), 0644)
// 	if err != nil {
// 		return "", err
// 	}

// 	uf, err := os.Open(userdataDestination)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer uf.Close()

// 	err = writer.AddFile(uf, "user-data")
// 	if err != nil {
// 		return "", err
// 	}

// 	networkConfigSource := path.Join(source, "user-data")
// 	networkConfig, err := os.ReadFile(networkConfigSource)
// 	if err != nil {
// 		return "", err
// 	}

// 	networkConfigDestination := path.Join(output, "network-config")
// 	err = os.WriteFile(networkConfigDestination, []byte(networkConfig), 0644)
// 	if err != nil {
// 		return "", err
// 	}

// 	nf, err := os.Open(networkConfigDestination)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer nf.Close()

// 	err = writer.AddFile(nf, "network-config")
// 	if err != nil {
// 		return "", err
// 	}

// 	destination := path.Join(output, "cloud-init.iso")
// 	of, err := os.OpenFile(destination, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer of.Close()

// 	err = writer.WriteTo(of, "cidata")
// 	if err != nil {
// 		return "", err
// 	}
// 	return destination, nil
// }

func createOverlayDisk() error {
	return nil
}
