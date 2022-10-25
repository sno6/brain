package brain

import (
	"os"

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
func (b *Brain) Write(s string) error {
	if s == "" {
		return nil
	}

	cell, err := b.buildCellFromData(s)
	if err != nil {
		return err
	}
	if err := b.writeCell(cell); err != nil {
		return err
	}
	return b.search.Index(cell.Identifier(), s)
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
func (b *Brain) List(query string, mode search.Mode) ([]*Cell, error) {
	ids, err := b.search.Query(query, mode)
	if err != nil {
		return nil, err
	}
	return b.readCells(ids)
}

// Delete removes a document from the index by its ID.
//
// The underlying data will remain in the data file, we are only
// removing the pointer to the data. We do this because document ids
// are offset references and deleting previous records would require
// a re-index of all records > n.
func (b *Brain) Delete(id string) error {
	return b.search.Delete(id)
}

func (b *Brain) buildCellFromData(data string) (*Cell, error) {
	offset, err := size(b.data)
	if err != nil {
		return nil, err
	}
	return NewCell(offset, data), nil
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
	return ParseCell(offset, string(buf))
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
