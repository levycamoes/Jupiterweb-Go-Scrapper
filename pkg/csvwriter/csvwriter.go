package csvwriter

import (
	"encoding/csv"
	"os"
)

func NewCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{"Ano", "Unidade", "Título", "Vertente", "Bolsas"})
	return writer, file, nil
}
