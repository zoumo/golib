package fileinfo

import "os"

func IsSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}
