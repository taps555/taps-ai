package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"fmt"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (fs *FileService) ProcessFile(content string) (map[string][]string, error) {
	// Membaca konten file CSV
	reader := csv.NewReader(strings.NewReader(content))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// Menyiapkan map untuk menyimpan data berdasarkan kolom
	table := make(map[string][]string)

	// Menyusun data sesuai dengan kolom pada dataset
	for _, record := range records[1:] { // Skipping header
		table["Date"] = append(table["Date"], record[0])
		table["Time"] = append(table["Time"], record[1])
		table["Appliance"] = append(table["Appliance"], record[2])
		table["Energy_Consumption"] = append(table["Energy_Consumption"], record[3])
		table["Room"] = append(table["Room"], record[4])
		table["Status"] = append(table["Status"], record[5])
	}

	return table, nil
}