package stringgrid

import "github.com/pkg/errors"

type ColumnGenerator interface {
	Generate(i uint) (string, error)
}

type Columns []ColumnGenerator

func (cs Columns) Generate(rows uint) ([][]string, error) {
	var data [][]string
	for i := uint(0); i < rows; i++ {
		row, err := generateRow(cs, i)
		if err != nil {
			return nil, errors.Wrapf(err, "generating row:%d", i)
		}
		data = append(data, row)
	}
	return data, nil
}

func generateRow(cs Columns, j uint) ([]string, error) {
	var row []string
	for _, c := range cs {
		cell, err := c.Generate(j)
		if err != nil {
			return nil, errors.Wrapf(err, "generating cell content at Column:%d", j)
		}
		row = append(row, cell)
	}
	return row, nil
}

type SimpleColumn func(i uint) string

func (sc SimpleColumn) Generate(i uint) (string, error) {
	return sc(i), nil
}

type ErrorColumn func(i uint) (string, error)

func (ec ErrorColumn) Generate(i uint) (string, error) {
	return ec(i)
}
