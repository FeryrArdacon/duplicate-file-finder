//go:build !windows

package filehandler

func isHiddenFile(filename string, directory string) (bool, error) {
	return filename[0] == '.', nil
}
