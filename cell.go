package brain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// A Cell is any individual idea / thought / note that is written
// to Brain. On disk, a cell is a byte array prepended with a UTC
// timestamp that starts at some byte offset in the .data file.
type Cell struct {
	offset int64
	ts     int64
	data   []byte
}

// NewCell returns a new cell with the given data.
func NewCell(offset int64, data []byte) *Cell {
	return &Cell{
		offset: offset,
		ts:     time.Now().UTC().Unix(),
		data:   data,
	}
}

func ParseCell(offset int64, data []byte) *Cell {
	return &Cell{
		offset: offset,
		data:   data,
		ts:     parseTimestamp(data),
	}
}

func (c *Cell) Identifier() string {
	return fmt.Sprintf("%d:%d", c.offset, len(c.Marshal()))
}

func (c *Cell) Marshal() []byte {
	return []byte(
		fmt.Sprintf("%d %s", c.ts, c.data),
	)
}

func (c *Cell) Data() []byte {
	return c.data
}

func (c *Cell) Timestamp() time.Time {
	return time.Unix(c.ts, 0)
}

func parseTimestamp(data []byte) int64 {
	n, _ := strconv.ParseInt(string(data[:10]), 10, 64)
	return n
}

func parseIdentifier(id string) (int64, int64, error) {
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return 0, 0, errors.New("invalid cell identifier")
	}

	offset, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	size, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return offset, size, nil
}
