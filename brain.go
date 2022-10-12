package brain

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/sno6/brain/search"
)

const (
	brainDir = ".brain/"
	dataFn   = ".data"
)

// Brain is the main driver for reading, writing, and querying
// your brain data file.
//
// Brain writes to the ~/.brain folder on your hard-drive, creating
// the following files:
//
// - .data which stores raw cell data.
// - .index.bleve/ which stores all index data used for querying.
//
type Brain struct {
	data   *os.File
	search *search.Search
}

// New initialises a new Brain.
func New() (*Brain, error) {
	data, dir, err := initBrain()
	if err != nil {
		return nil, err
	}
	s, err := search.New(dir)
	if err != nil {
		return nil, err
	}
	return &Brain{
		data:   data,
		search: s,
	}, nil
}

// Write spawns an editor to capture user input, and pipes the bytes to
// the .data file. It then indexes the content for future queries.
func (b *Brain) Write() error {
	data, err := captureEditor()
	if err != nil {
		return err
	}

	cell, err := b.buildCellFromData(data)
	if err != nil {
		return err
	}
	if err := b.writeCell(cell); err != nil {
		return err
	}

	return b.search.Index(cell.Identifier(), string(data))
}

// Read reads a cell in .data by a given identifier.
func (b *Brain) Read(id string) (*Cell, error) {
	offset, sz, err := parseIdentifier(id)
	if err != nil {
		return nil, err
	}
	return b.readCell(offset, sz)
}

// List searches for cells within .data by checking the index against
// a given query and returns any cells that match.
func (b *Brain) List(query string) ([]*Cell, error) {
	ids, err := b.search.Query(query)
	if err != nil {
		return nil, err
	}
	return b.readCells(ids)
}

func (b *Brain) buildCellFromData(data []byte) (*Cell, error) {
	offset, err := size(b.data)
	if err != nil {
		return nil, err
	}
	return NewCellFromData(offset, data), nil
}

func (b *Brain) writeCell(cell *Cell) error {
	_, err := b.data.Write(cell.Marshal())
	return err
}

func (b *Brain) readCell(offset, size int64) (*Cell, error) {
	buf := make([]byte, size)
	_, err := b.data.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}
	return NewCellFromData(offset, buf), nil
}

func (b *Brain) readCells(ids []string) ([]*Cell, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	cells := make([]*Cell, len(ids))
	for i, id := range ids {
		cell, err := b.Read(id)
		if err != nil {
			return nil, err
		}
		cells[i] = cell
	}

	return cells, nil
}

func captureEditor() ([]byte, error) {
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
