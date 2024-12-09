package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"fmt"
	"log"
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

    // Log untuk memeriksa isi CSV
    log.Printf("CSV Records: %+v", records)

    // Pastikan ada baris setelah header
    if len(records) <= 1 {
        return nil, fmt.Errorf("no data found in CSV after the header")
    }

    // Menyiapkan map untuk menyimpan data berdasarkan kolom
    table := make(map[string][]string)

    // Memeriksa baris pertama sebagai header dan melanjutkan ke baris berikutnya
    for _, record := range records[1:] { // Skipping header
        if len(record) < 6 { // Memeriksa apakah baris memiliki cukup kolom
            log.Printf("Skipping record (too few columns): %v", record)
            continue
        }

        // Memeriksa apakah baris kosong atau memiliki data yang tidak valid
        if record[0] == "" || record[1] == "" || record[2] == "" || record[3] == "" || record[4] == "" || record[5] == "" {
            log.Printf("Skipping incomplete record: %v", record)
            continue
        }

        log.Println(record)

        // Menyusun data sesuai dengan kolom pada dataset
        table["Date"] = append(table["Date"], record[0])
        table["Time"] = append(table["Time"], record[1])
        table["Appliance"] = append(table["Appliance"], record[2])
        table["Energy_Consumption"] = append(table["Energy_Consumption"], record[3])
        table["Room"] = append(table["Room"], record[4])
        table["Status"] = append(table["Status"], record[5])
    }

	log.Println("File content:", content)

    // Pastikan ada data yang valid dalam tabel
    if len(table["Date"]) == 0 {
        return nil, fmt.Errorf("no valid records found in CSV")
    }

    return table, nil
}

