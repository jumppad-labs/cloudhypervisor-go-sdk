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

	// figure out what the minimum size for an ISO9660 disk is.
	// 10MB works but we could perhaps use less.
	var diskSize int64 = 10 * 1024 * 1024
	if size < diskSize {
		size = diskSize
	}

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
		return fmt.Errorf("not an iso9660 filesystem")
	}

	err = iso.Finalize(iso9660.FinalizeOptions{})
	if err != nil {
		return err
	}

	return nil
}

func createOverlayDisk() error {
	return nil
}
