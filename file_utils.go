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

// Create FS resources if they don't yet exist.
func initBrain() (*os.File, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}

	dir := path.Join(home, brainDir)
	if _, err = os.Stat(dir); err != nil && os.IsNotExist(err) {
		// First time running.. init resources.
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return nil, "", err
		}
		if _, err := os.Create(path.Join(dir, dataFn)); err != nil {
			return nil, "", err
		}
	}

	f, err := os.OpenFile(path.Join(dir, dataFn), os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return nil, "", err
	}
	return f, dir, nil
}
