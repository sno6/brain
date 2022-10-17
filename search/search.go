package search

import (
	"path"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

const indexFn = ".index.bleve"

// Mode is an enum that defines the type of search to run against the index.
type Mode uint8

const (
	Match Mode = iota
	Fuzzy
	Wildcard
)

// Search is responsible for creating and operating a bleve index.
type Search struct {
	index bleve.Index
}

// New initialises Search with an open index ready for querying.
//
// If it's the first time this has been called it will initialise
// a new folder for the index under the given directory.
func New(dir string) (*Search, error) {
	index, err := openIndexOrInit(dir)
	if err != nil {
		return nil, err
	}
	return &Search{index: index}, nil
}

// Index indexes the data for a given id.
//
// The identifier we use here is a combination of the byte offset
// and size of the data in bytes which we can use to quickly read the
// data from disk.
func (s *Search) Index(id string, data string) error {
	return s.index.Index(id, data)
}

// Query runs a match query on the index and returns the document
// ids if there are any matches
func (s *Search) Query(qs string, mode Mode) ([]string, error) {
	var q query.Query
	switch mode {
	case Match:
		q = bleve.NewMatchQuery(qs)
	case Wildcard:
		q = bleve.NewWildcardQuery(qs)
	case Fuzzy:
		q = bleve.NewFuzzyQuery(qs)
	default:
		q = bleve.NewMatchQuery(qs)
	}

	r := bleve.NewSearchRequest(q)

	// How many do we want to return in a single request?
	r.Size = 100

	res, err := s.index.Search(r)
	if err != nil {
		return nil, err
	}

	if res == nil || len(res.Hits) == 0 {
		return nil, nil
	}

	ids := make([]string, len(res.Hits))
	for i, h := range res.Hits {
		ids[i] = h.ID
	}

	return ids, nil
}

func openIndexOrInit(dir string) (bleve.Index, error) {
	fullPath := path.Join(dir, indexFn)

	index, err := bleve.Open(fullPath)
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			return nil, err
		}

		// Initialise a new index with default mappings.
		return bleve.New(
			fullPath,
			bleve.NewIndexMapping(),
		)
	}

	return index, nil
}
