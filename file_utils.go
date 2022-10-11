package brain

import (
	"os"
	"path"
)

// size returns the size in bytes of the given file.
func size(f *os.File) (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func initFS(dir string) error {
	_, err := os.Stat(dir)
	if err == nil || os.IsExist(err) {
		return nil
	}
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}
	_, err = os.OpenFile(path.Join(dir, ".data"), os.O_WRONLY|os.O_CREATE, 0755)
	return err
}
