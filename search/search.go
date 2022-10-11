package search

import (
	"path"

	"github.com/blevesearch/bleve/v2"
)

const indexFn = ".index.bleve"

// Search is responsible for creating and operating indexes.
type Search struct {
	dir   string
	index bleve.Index
}

// New initialises Search with an open index reading for querying.
func New(dir string) (*Search, error) {
	index, err := openIndexOrInit(dir)
	if err != nil {
		return nil, err
	}
	return &Search{
		dir:   dir,
		index: index,
	}, nil
}

// Index indexes the data for a given id.
//
// The identifier we use here is a combination of the byte offset
// and size of the data in bytes.
func (s *Search) Index(id string, data string) error {
	return s.index.Index(id, data)
}

func (s *Search) Query(query string) ([]string, error) {
	q := bleve.NewMatchQuery(query)
	r := bleve.NewSearchRequest(q)
	res, err := s.index.Search(r)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, h := range res.Hits {
		ids = append(ids, h.ID)
	}

	return ids, nil
}

func openIndexOrInit(dir string) (bleve.Index, error) {
	p := path.Join(dir, indexFn)

	index, err := bleve.Open(p)
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			return nil, err
		}

		// Initialise a new index with default mappings.
		return bleve.New(
			p,
			bleve.NewIndexMapping(),
		)
	}

	return index, nil
}
