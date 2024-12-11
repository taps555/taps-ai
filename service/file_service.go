package service

import (
	"a21hc3NpZ25tZW50/model"
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"strconv"
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

func (fs *FileService) ProcessTableData(table map[string][]string, question string) (model.TapasResponse, error) {
    var answer string
	var coordinates [][]int // Untuk menyimpan posisi koordinat jika diperlukan
	var cells []string      // Menyimpan sel terkait jika diperlukan
	var aggregator string   // Menyimpan jenis agregasi jika diperlukan

	applianceEnergy := make(map[string]float64)
	applianceCount := make(map[string]int)

	log.Println("Data perangkat yang terdeteksi:")
	for i, appliance := range table["Appliance"] {
		log.Printf("Perangkat %d: %s\n", i+1, appliance)
	}

	questionLower := strings.ToLower(question)

	// Menangani pertanyaan tentang total energi untuk perangkat tertentu
	if strings.Contains(questionLower, "total energi") {
		var applianceQuery string
		if strings.Contains(questionLower, "ac") {
			applianceQuery = "AC"
		} else if strings.Contains(questionLower, "tv") {
			applianceQuery = "TV"
		} else if strings.Contains(questionLower, "evcar") {
			applianceQuery = "EVCar"
		} else if strings.Contains(questionLower, "refrigerator") {
			applianceQuery = "Refrigerator"
		} else if strings.Contains(questionLower, "fridge") { // Alias untuk Refrigerator
			applianceQuery = "Refrigerator"
		} else if strings.Contains(questionLower, "televisi") { // Alias untuk TV
			applianceQuery = "TV"
		}

		log.Printf("Mencari perangkat: %s\n", applianceQuery)

		if applianceQuery != "" {
			found := false
			for i := 0; i < len(table["Appliance"]); i++ {
				appliance := table["Appliance"][i]
				applianceLower := strings.ToLower(appliance)

				if applianceLower == strings.ToLower(applianceQuery) { // Pastikan perbandingan case-insensitive
					found = true
					if _, exists := applianceEnergy[appliance]; !exists {
						energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
						if err != nil {
							return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
						}
						applianceEnergy[appliance] += energyConsumption
						applianceCount[appliance]++
					}
				}
			}

			if found {
				answer = fmt.Sprintf("%s: %.2f kWh.", applianceQuery, applianceEnergy[applianceQuery])
			} else {
				answer = fmt.Sprintf("Perangkat %s tidak ditemukan dalam data.", applianceQuery)
			}
		} else {
			answer = "Perangkat yang diminta tidak dikenali dalam pertanyaan."
		}
	}

	// Menangani pertanyaan tentang total energi semua perangkat
	if strings.Contains(questionLower, "total energi semua perangkat") {
		totalEnergy := 0.0
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			totalEnergy += energyConsumption
		}
		answer = fmt.Sprintf("Total energi dari semua perangkat: %.2f kWh.", totalEnergy)
	}

	// Menangani pertanyaan tentang konsumsi energi per hari
	if strings.Contains(questionLower, "konsumsi energi per hari") {
		dailyEnergy := make(map[string]float64)
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			appliance := table["Appliance"][i]
			dailyEnergy[appliance] += energyConsumption
		}
		var dailyReport []string
		for appliance, energy := range dailyEnergy {
			dailyReport = append(dailyReport, fmt.Sprintf("%s: %.2f kWh", appliance, energy))
		}
		answer = "Konsumsi Energi Per Hari: " + strings.Join(dailyReport, ", ")
	}

	// Menangani pertanyaan tentang konsumsi energi per minggu
	if strings.Contains(questionLower, "konsumsi energi per minggu") {
		weeklyEnergy := make(map[string]float64)
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			appliance := table["Appliance"][i]
			weeklyEnergy[appliance] += energyConsumption * 7 // Anggap konsumsi harian sama setiap hari dalam seminggu
		}
		var weeklyReport []string
		for appliance, energy := range weeklyEnergy {
			weeklyReport = append(weeklyReport, fmt.Sprintf("%s: %.2f kWh", appliance, energy))
		}
		answer = "Konsumsi Energi Per Minggu: " + strings.Join(weeklyReport, ", ")
	}

	// Menangani pertanyaan tentang konsumsi energi per bulan
	if strings.Contains(questionLower, "konsumsi energi per bulan") {
		monthlyEnergy := make(map[string]float64)
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			appliance := table["Appliance"][i]
			monthlyEnergy[appliance] += energyConsumption * 30 // Anggap konsumsi harian sama setiap hari dalam sebulan
		}
		var monthlyReport []string
		for appliance, energy := range monthlyEnergy {
			monthlyReport = append(monthlyReport, fmt.Sprintf("%s: %.2f kWh", appliance, energy))
		}
		answer = "Konsumsi Energi Per Bulan: " + strings.Join(monthlyReport, ", ")
	}

	// Menangani pertanyaan tentang perangkat dengan konsumsi energi tertinggi
	if strings.Contains(questionLower, "perangkat dengan konsumsi energi tertinggi") {
		highestEnergy := 0.0
		applianceWithHighestEnergy := ""
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			appliance := table["Appliance"][i]
			if energyConsumption > highestEnergy {
				highestEnergy = energyConsumption
				applianceWithHighestEnergy = appliance
			}
		}
		answer = fmt.Sprintf("Perangkat dengan konsumsi energi tertinggi adalah %s dengan konsumsi energi %.2f kWh.", applianceWithHighestEnergy, highestEnergy)
	}

	// Menangani pertanyaan tentang perangkat dengan konsumsi energi terendah
	if strings.Contains(questionLower, "perangkat dengan konsumsi energi terendah") {
		lowestEnergy := math.MaxFloat64
		applianceWithLowestEnergy := ""
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			appliance := table["Appliance"][i]
			if energyConsumption < lowestEnergy {
				lowestEnergy = energyConsumption
				applianceWithLowestEnergy = appliance
			}
		}
		answer = fmt.Sprintf("Perangkat dengan konsumsi energi terendah adalah %s dengan konsumsi energi %.2f kWh.", applianceWithLowestEnergy, lowestEnergy)
	}

	// Menangani pertanyaan tentang perangkat yang tidak aktif
	if strings.Contains(questionLower, "perangkat tidak aktif") {
		inactiveAppliances := []string{}
		for i := 0; i < len(table["Appliance"]); i++ {
			if strings.ToLower(table["Status"][i]) == "off" {
				appliance := table["Appliance"][i]
				inactiveAppliances = append(inactiveAppliances, appliance)
			}
		}
		if len(inactiveAppliances) > 0 {
			answer = "Perangkat yang tidak aktif: " + strings.Join(inactiveAppliances, ", ")
		} else {
			answer = "Tidak ada perangkat yang tidak aktif."
		}
	}

	// Menangani pertanyaan tentang total penghematan energi
	if strings.Contains(questionLower, "total penghematan energi") {
		totalSavings := 0.0
		for i := 0; i < len(table["Appliance"]); i++ {
			if strings.ToLower(table["Status"][i]) == "off" {
				energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
				if err != nil {
					return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
				}
				totalSavings += energyConsumption
			}
		}
		answer = fmt.Sprintf("Total penghematan energi dengan perangkat yang dimatikan: %.2f kWh.", totalSavings)
	}

	// Menangani pertanyaan tentang perangkat di lokasi tertentu
	if strings.Contains(questionLower, "perangkat di") {
		var roomQuery string
		if strings.Contains(questionLower, "dapur") {
			roomQuery = "Kitchen"
		} else if strings.Contains(questionLower, "ruang tamu") {
			roomQuery = "Living Room"
		} else if strings.Contains(questionLower, "kamar tidur") {
			roomQuery = "Bedroom"
		}

		log.Printf("Mencari perangkat di lokasi: %s\n", roomQuery)

		if roomQuery != "" {
			appliancesInRoom := []string{}
			for i := 0; i < len(table["Room"]); i++ {
				if strings.ToLower(table["Room"][i]) == strings.ToLower(roomQuery) {
					appliancesInRoom = append(appliancesInRoom, table["Appliance"][i])
				}
			}
			if len(appliancesInRoom) > 0 {
				answer = "Perangkat di " + roomQuery + ": " + strings.Join(appliancesInRoom, ", ")
			} else {
				answer = "Tidak ada perangkat di " + roomQuery + "."
			}
		} else {
			answer = "Lokasi yang diminta tidak dikenali dalam pertanyaan."
		}
	}

        // Menangani pertanyaan tentang perangkat dengan konsumsi energi rata-rata tertinggi
    if strings.Contains(questionLower, "perangkat dengan konsumsi energi rata-rata tertinggi") {
        applianceEnergy := make(map[string]float64)
        applianceCount := make(map[string]int)
        
        for i := 0; i < len(table["Appliance"]); i++ {
            appliance := table["Appliance"][i]
            energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
            }
            applianceEnergy[appliance] += energyConsumption
            applianceCount[appliance]++
        }
        
        highestAvgEnergy := 0.0
        applianceWithHighestAvgEnergy := ""
        for appliance, totalEnergy := range applianceEnergy {
            averageEnergy := totalEnergy / float64(applianceCount[appliance])
            if averageEnergy > highestAvgEnergy {
                highestAvgEnergy = averageEnergy
                applianceWithHighestAvgEnergy = appliance
            }
        }
        
        answer = fmt.Sprintf("Perangkat dengan konsumsi energi rata-rata tertinggi adalah %s dengan rata-rata konsumsi energi %.2f kWh.", applianceWithHighestAvgEnergy, highestAvgEnergy)
    }

    // Menangani pertanyaan tentang perangkat dengan konsumsi energi lebih tinggi dari nilai tertentu
    if strings.Contains(questionLower, "perangkat dengan konsumsi energi lebih tinggi dari") {
        threshold := 5.0 // Misalnya, ambang batas konsumsi energi
        highEnergyAppliances := []string{}
        for i := 0; i < len(table["Energy_Consumption"]); i++ {
            energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
            }
            if energyConsumption > threshold {
                appliance := table["Appliance"][i]
                highEnergyAppliances = append(highEnergyAppliances, appliance)
            }
        }
        if len(highEnergyAppliances) > 0 {
            answer = "Perangkat dengan konsumsi energi lebih tinggi dari " + fmt.Sprintf("%.2f", threshold) + " kWh: " + strings.Join(highEnergyAppliances, ", ")
        } else {
            answer = fmt.Sprintf("Tidak ada perangkat dengan konsumsi energi lebih tinggi dari %.2f kWh.", threshold)
        }
    }

    // Menangani pertanyaan tentang perangkat dengan konsumsi energi lebih rendah dari nilai tertentu
    if strings.Contains(questionLower, "perangkat dengan konsumsi energi lebih rendah dari") {
        threshold := 1.0 // Misalnya, ambang batas konsumsi energi
        lowEnergyAppliances := []string{}
        for i := 0; i < len(table["Energy_Consumption"]); i++ {
            energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
            }
            if energyConsumption < threshold {
                appliance := table["Appliance"][i]
                lowEnergyAppliances = append(lowEnergyAppliances, appliance)
            }
        }
        if len(lowEnergyAppliances) > 0 {
            answer = "Perangkat dengan konsumsi energi lebih rendah dari " + fmt.Sprintf("%.2f", threshold) + " kWh: " + strings.Join(lowEnergyAppliances, ", ")
        } else {
            answer = fmt.Sprintf("Tidak ada perangkat dengan konsumsi energi lebih rendah dari %.2f kWh.", threshold)
        }
    }

    // Menangani pertanyaan tentang jumlah total perangkat
    if strings.Contains(questionLower, "jumlah total perangkat") {
        totalCount := len(table["Appliance"])
        answer = fmt.Sprintf("Jumlah total perangkat adalah %d.", totalCount)
    }

    // Menangani pertanyaan tentang perangkat dengan status tertentu (On atau Off)
    if strings.Contains(questionLower, "perangkat dengan status on") {
        activeAppliances := []string{}
        for i := 0; i < len(table["Status"]); i++ {
            if strings.ToLower(table["Status"][i]) == "on" {
                appliance := table["Appliance"][i]
                activeAppliances = append(activeAppliances, appliance)
            }
        }
        if len(activeAppliances) > 0 {
            answer = "Perangkat yang aktif: " + strings.Join(activeAppliances, ", ")
        } else {
            answer = "Tidak ada perangkat yang aktif."
        }
    }

    if strings.Contains(questionLower, "perangkat dengan status off") {
        inactiveAppliances := []string{}
        for i := 0; i < len(table["Status"]); i++ {
            if strings.ToLower(table["Status"][i]) == "off" {
                appliance := table["Appliance"][i]
                inactiveAppliances = append(inactiveAppliances, appliance)
            }
        }
        if len(inactiveAppliances) > 0 {
            answer = "Perangkat yang mati: " + strings.Join(inactiveAppliances, "," + "\n")
        } else {
            answer = "Tidak ada perangkat yang mati."
        }
    }

    // Menangani pertanyaan tentang total biaya energi
    if strings.Contains(questionLower, "total biaya energi") {
        ratePerKWh := 1500.0 // Tarif energi per kWh (dalam Rupiah)
        totalCost := 0.0
        for i := 0; i < len(table["Energy_Consumption"]); i++ {
            energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
            }
            totalCost += energyConsumption * ratePerKWh
        }
        answer = fmt.Sprintf("Total biaya energi untuk semua perangkat adalah Rp %.2f.", totalCost)
    }
	return model.TapasResponse{
		Answer: answer,
		Coordinates: coordinates,  // Tambahkan jika perlu koordinat
		Cells:       cells,        // Tambahkan jika perlu sel
		Aggregator:  aggregator,    // Tambahkan jika perlu agregator
	}, nil

}

