package brain

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/sno6/brain/search"
)

type Brain struct {
	dir    string
	search *search.Search
}

func New() (*Brain, error) {
	return newBrain()
}

func (b *Brain) Write() error {
	// Spawn an editor to collect cell data from the user.
	data, err := handleEditor()
	if err != nil {
		return err
	}

	f, err := b.openData(os.O_WRONLY | os.O_APPEND)
	if err != nil {
		return err
	}
	defer f.Close()

	// The current size of the file before writing data is the offset
	// for the next cell. This offset is used as the identifier for the index.
	offset, err := size(f)
	if err != nil {
		return err
	}

	cell := NewCellFromData(offset, data)
	if _, err := f.Write(cell.Marshal()); err != nil {
		return err
	}

	// Now that we have written the cell to disk we can index
	// the contents, using the offset and size of the cell as the
	// identifier.
	return b.search.Index(cell.Identifier(), string(data))
}

func (b *Brain) Read(id string) (string, error) {
	offset, size, err := parseCellID(id)
	if err != nil {
		return "", err
	}
	return b.readCellData(offset, size)
}

func (b *Brain) List(query string) ([]string, error) {
	return b.search.Query(query)
}

func (b *Brain) readCellData(offset, size int) (string, error) {
	f, err := b.openData(os.O_RDONLY)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, size)
	if _, err := f.ReadAt(buf, int64(offset)); err != nil {
		return "", err
	}

	return string(buf), nil
}

func (b *Brain) openData(flags int) (*os.File, error) {
	return os.OpenFile(path.Join(b.dir, ".data"), flags, 0755)
}

func newBrain() (*Brain, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := path.Join(home, ".brain/")
	if err := initFS(dir); err != nil {
		return nil, err
	}
	s, err := search.New(dir)
	if err != nil {
		return nil, err
	}
	return &Brain{
		dir:    dir,
		search: s,
	}, nil
}

func handleEditor() ([]byte, error) {
	tmp, err := os.CreateTemp("/tmp", "brain")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			log.Printf("Unable to remove temp file: %v\n", err)
		}
	}()

	cmd := exec.Command("vim", tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return ioutil.ReadAll(tmp)
}
