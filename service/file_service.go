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

		var energyReports []string
		dailyEnergy := make(map[string]float64)
		weeklyEnergy := make(map[string]float64)
		monthlyEnergy := make(map[string]float64)

		for i := 0; i < len(table["Appliance"]); i++ {
			appliance := table["Appliance"][i]
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				log.Printf("Error parsing energy for appliance %s: %v", appliance, err)
				continue
			}

			dailyEnergy[appliance] += energyConsumption
			weeklyEnergy[appliance] += energyConsumption * 7
			monthlyEnergy[appliance] += energyConsumption * 30
		}

		for appliance, daily := range dailyEnergy {
			weekly := weeklyEnergy[appliance]
			monthly := monthlyEnergy[appliance]
			energyReports = append(energyReports, fmt.Sprintf("%s: Daily %.2f kWh, Weekly %.2f kWh, Monthly %.2f kWh", appliance, daily, weekly, monthly))
		}

		answer = "Energy Consumption Reports:\n" + strings.Join(energyReports, "\n")

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
		energyByDevice := make(map[string]float64)
		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumptionStr := table["Energy_Consumption"][i]
			if energyConsumptionStr == "" {
				continue
			}
	
			energyConsumption, err := strconv.ParseFloat(energyConsumptionStr, 64)
			if err != nil {
				log.Printf("Error parsing energy for %s: %v", table["Appliance"][i], err)
				continue
			}
	
			appliance := table["Appliance"][i]
			totalEnergy += energyConsumption
			energyByDevice[appliance] += energyConsumption
		}
	
		if totalEnergy > 0 {
			var report []string
			for device, energy := range energyByDevice {
				percentage := (energy / totalEnergy) * 100
				daily := energy / 7
				monthly := energy * 4
				report = append(report, fmt.Sprintf("%s: %.2f kWh (%.2f%%) - Daily: %.2f kWh, Weekly: %.2f kWh, Monthly: %.2f kWh",
					device, energy, percentage, daily, energy, monthly))
			}
			answer = "Energy Breakdown: " + strings.Join(report, "; ")
		} else {
			answer = "No energy data available for devices."
		}
	}
	

	// Menangani pertanyaan tentang konsumsi energi per hari
	if strings.Contains(questionLower, "konsumsi energi per hari") {
        dailyEnergy := make(map[string]float64)
        for i := 0; i < len(table["Appliance"]); i++ {
            energyStr := table["Energy_Consumption"][i]
            
            // Cek apakah energyStr kosong atau tidak valid
            if energyStr == "" {
                log.Printf("Warning: Energy consumption is empty for appliance %s\n", table["Appliance"][i])
                continue // Skip jika tidak ada nilai energi
            }

            energyConsumption, err := strconv.ParseFloat(energyStr, 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption for %s: %v", table["Appliance"][i], err)
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
            energyStr := table["Energy_Consumption"][i]
            
            // Cek apakah energyStr kosong atau tidak valid
            if energyStr == "" {
                log.Printf("Warning: Energy consumption is empty for appliance %s\n", table["Appliance"][i])
                continue // Skip jika tidak ada nilai energi
            }

            energyConsumption, err := strconv.ParseFloat(energyStr, 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption for %s: %v", table["Appliance"][i], err)
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
            energyStr := table["Energy_Consumption"][i]
            
            // Cek apakah energyStr kosong atau tidak valid
            if energyStr == "" {
                log.Printf("Warning: Energy consumption is empty for appliance %s\n", table["Appliance"][i])
                continue // Skip jika tidak ada nilai energi
            }

            energyConsumption, err := strconv.ParseFloat(energyStr, 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption for %s: %v", table["Appliance"][i], err)
            }

            appliance := table["Appliance"][i]
            monthlyEnergy[appliance] += energyConsumption * 30 // Anggap konsumsi harian sama setiap hari dalam sebulan (30 hari)
        }

        var monthlyReport []string
        for appliance, energy := range monthlyEnergy {
            monthlyReport = append(monthlyReport, fmt.Sprintf("%s: %.2f kWh", appliance, energy))
        }
        answer = "Konsumsi Energi Per Bulan: " + strings.Join(monthlyReport, ", ")
    }

	// Menangani pertanyaan tentang perangkat dengan konsumsi energi tertinggi
	// Handle devices with the highest or lowest energy consumption
	if strings.Contains(questionLower, "konsumsi energi tertinggi") || strings.Contains(questionLower, "konsumsi energi terendah") {
		isHighest := strings.Contains(questionLower, "tertinggi")
		targetEnergy := 0.0
		if !isHighest {
			targetEnergy = math.MaxFloat64
		}
		targetAppliance := ""

		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			appliance := table["Appliance"][i]
			if (isHighest && energyConsumption > targetEnergy) || (!isHighest && energyConsumption < targetEnergy) {
				targetEnergy = energyConsumption
				targetAppliance = appliance
			}
		}

		comparison := "tertinggi"
		if !isHighest {
			comparison = "terendah"
		}
		answer = fmt.Sprintf("Perangkat dengan konsumsi energi %s adalah %s dengan konsumsi energi %.2f kWh.", comparison, targetAppliance, targetEnergy)
	}

// Handle devices with energy above or below a certain threshold
	if strings.Contains(questionLower, "lebih tinggi dari") || strings.Contains(questionLower, "lebih rendah dari") {
		var threshold float64
		if strings.Contains(questionLower, "lebih tinggi dari") {
			fmt.Sscanf(questionLower, "lebih tinggi dari %f", &threshold)
		} else if strings.Contains(questionLower, "lebih rendah dari") {
			fmt.Sscanf(questionLower, "lebih rendah dari %f", &threshold)
		}

		highEnergy := strings.Contains(questionLower, "lebih tinggi dari")
		appliances := []string{}
		for i := 0; i < len(table["Energy_Consumption"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
			}
			if (highEnergy && energyConsumption > threshold) || (!highEnergy && energyConsumption < threshold) {
				appliance := table["Appliance"][i]
				appliances = append(appliances, appliance)
			}
		}

		if len(appliances) > 0 {
			comparison := "lebih tinggi"
			if !highEnergy {
				comparison = "lebih rendah"
			}
			answer = fmt.Sprintf("Perangkat dengan konsumsi energi %s dari %.2f kWh: %s", comparison, threshold, strings.Join(appliances, ", "))
		} else {
			answer = fmt.Sprintf("Tidak ada perangkat dengan konsumsi energi sesuai kriteria (%.2f kWh).", threshold)
		}
	}

// Handle devices in a specific location
	if strings.Contains(questionLower, "perangkat di") {
		var roomQuery string
		if strings.Contains(questionLower, "dapur") {
			roomQuery = "Kitchen"
		} else if strings.Contains(questionLower, "ruang tamu") {
			roomQuery = "Living Room"
		} else if strings.Contains(questionLower, "kamar tidur") {
			roomQuery = "Bedroom"
		} else {
			answer = "Lokasi yang diminta tidak dikenali dalam pertanyaan."
			return model.TapasResponse{}, nil
		}

		appliancesInRoom := []string{}
		for i := 0; i < len(table["Room"]); i++ {
			if strings.EqualFold(table["Room"][i], roomQuery) {
				appliancesInRoom = append(appliancesInRoom, table["Appliance"][i])
			}
		}

		if len(appliancesInRoom) > 0 {
			answer = fmt.Sprintf("Perangkat di %s: %s", roomQuery, strings.Join(appliancesInRoom, ", "))
		} else {
			answer = fmt.Sprintf("Tidak ada perangkat di %s.", roomQuery)
		}
	}

	// Generate energy consumption reports (daily, weekly, monthly)



		// Menangani pertanyaan tentang perangkat dengan status tertentu (On atau Off)
		// Menangani pertanyaan tentang perangkat dengan status On
	if strings.Contains(questionLower, "perangkat dengan status on") {
		activeAppliances := make(map[string]bool) // Use a map to ensure uniqueness
		for i := 0; i < len(table["Status"]); i++ {
			if strings.ToLower(table["Status"][i]) == "on" {
				appliance := table["Appliance"][i]
				activeAppliances[appliance] = true
			}
		}

		if len(activeAppliances) > 0 {
			var uniqueActiveAppliances []string
			for appliance := range activeAppliances {
				uniqueActiveAppliances = append(uniqueActiveAppliances, appliance)
			}
			answer = "Perangkat yang aktif: " + strings.Join(uniqueActiveAppliances, ", ")
		} else {
			answer = "Tidak ada perangkat yang aktif."
		}
	}

	// Menangani pertanyaan tentang perangkat dengan status Off
	if strings.Contains(questionLower, "perangkat dengan status off") {
		inactiveAppliances := make(map[string]bool) // Use a map to ensure uniqueness
		for i := 0; i < len(table["Status"]); i++ {
			if strings.ToLower(table["Status"][i]) == "off" {
				appliance := table["Appliance"][i]
				inactiveAppliances[appliance] = true
			}
		}

		if len(inactiveAppliances) > 0 {
			var uniqueInactiveAppliances []string
			for appliance := range inactiveAppliances {
				uniqueInactiveAppliances = append(uniqueInactiveAppliances, appliance)
			}
			answer = "Perangkat yang mati: " + strings.Join(uniqueInactiveAppliances, ", ")
		} else {
			answer = "Tidak ada perangkat yang mati."
		}
	}

	// Menangani pertanyaan tentang perangkat tidak aktif
	if strings.Contains(questionLower, "perangkat tidak aktif") {
		inactiveAppliances := make(map[string]bool) // Use a map to ensure uniqueness
		for i := 0; i < len(table["Status"]); i++ {
			if strings.ToLower(table["Status"][i]) == "off" {
				appliance := table["Appliance"][i]
				inactiveAppliances[appliance] = true
			}
		}

		if len(inactiveAppliances) > 0 {
			var uniqueInactiveAppliances []string
			for appliance := range inactiveAppliances {
				uniqueInactiveAppliances = append(uniqueInactiveAppliances, appliance)
			}
			answer = "Perangkat yang tidak aktif: " + strings.Join(uniqueInactiveAppliances, ", ")
		} else {
			answer = "Tidak ada perangkat yang tidak aktif."
		}
	}

	
	// Menangani pertanyaan tentang total penghematan energi
	if strings.Contains(questionLower, "total penghematan energi") {
		totalSavings := 0.0
		for i := 0; i < len(table["Appliance"]); i++ {
			if strings.ToLower(table["Status"][i]) == "Off" {
				energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
				if err != nil {
					return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
				}
				totalSavings += energyConsumption
			}
		}
		answer = fmt.Sprintf("Total penghematan energi dengan perangkat yang dimatikan: %.2f kWh.", totalSavings)
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





