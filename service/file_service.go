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
    reader := csv.NewReader(strings.NewReader(content))
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("error reading CSV: %v", err)
    }

    log.Printf("CSV Records: %+v", records)

    if len(records) <= 1 {
        return nil, fmt.Errorf("no data found in CSV after the header")
    }

    table := make(map[string][]string)

    for _, record := range records[1:] { 
        if len(record) < 6 { 
            log.Printf("Skipping record (too few columns): %v", record)
            continue
        }

        if record[0] == "" || record[1] == "" || record[2] == "" || record[3] == "" || record[4] == "" || record[5] == "" {
            log.Printf("Skipping incomplete record: %v", record)
            continue
        }

        log.Println(record)

        table["Date"] = append(table["Date"], record[0])
        table["Time"] = append(table["Time"], record[1])
        table["Appliance"] = append(table["Appliance"], record[2])
        table["Energy_Consumption"] = append(table["Energy_Consumption"], record[3])
        table["Room"] = append(table["Room"], record[4])
        table["Status"] = append(table["Status"], record[5])
    }

	log.Println("File content:", content)

 
    if len(table["Date"]) == 0 {
        return nil, fmt.Errorf("no valid records found in CSV")
    }

    return table, nil
}

func (fs *FileService) ProcessTableData(table map[string][]string, question string) (model.TapasResponse, error) {
    var answer string
    var coordinates [][]int 
    var cells []string    
    var aggregator string  

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


		questionLower := strings.ToLower(question)

	
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
		} else if strings.Contains(questionLower, "fridge") { 
			applianceQuery = "Refrigerator"
		} else if strings.Contains(questionLower, "televisi") {
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

	if strings.Contains(questionLower, "energi semua perangkat") || strings.Contains(questionLower, "konsumsi energi semua perangkat") {
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
	

	if strings.Contains(questionLower, "konsumsi energi per hari") {
        dailyEnergy := make(map[string]float64)
        for i := 0; i < len(table["Appliance"]); i++ {
            energyStr := table["Energy_Consumption"][i]
            
            if energyStr == "" {
                log.Printf("Warning: Energy consumption is empty for appliance %s\n", table["Appliance"][i])
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

    if strings.Contains(questionLower, "konsumsi energi per minggu") {
        weeklyEnergy := make(map[string]float64)
        for i := 0; i < len(table["Appliance"]); i++ {
            energyStr := table["Energy_Consumption"][i]
            
        
            if energyStr == "" {
                log.Printf("Warning: Energy consumption is empty for appliance %s\n", table["Appliance"][i])
         
            }

            energyConsumption, err := strconv.ParseFloat(energyStr, 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption for %s: %v", table["Appliance"][i], err)
            }
            
            appliance := table["Appliance"][i]
            weeklyEnergy[appliance] += energyConsumption * 7 
        }

        var weeklyReport []string
        for appliance, energy := range weeklyEnergy {
            weeklyReport = append(weeklyReport, fmt.Sprintf("%s: %.2f kWh", appliance, energy))
        }
        answer = "Konsumsi Energi Per Minggu: " + strings.Join(weeklyReport, ", ")
    }


    if strings.Contains(questionLower, "konsumsi energi per bulan") {
        monthlyEnergy := make(map[string]float64)
        for i := 0; i < len(table["Appliance"]); i++ {
            energyStr := table["Energy_Consumption"][i]
            
           
            if energyStr == "" {
                log.Printf("Warning: Energy consumption is empty for appliance %s\n", table["Appliance"][i])
              
            }

            energyConsumption, err := strconv.ParseFloat(energyStr, 64)
            if err != nil {
                return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption for %s: %v", table["Appliance"][i], err)
            }

            appliance := table["Appliance"][i]
            monthlyEnergy[appliance] += energyConsumption * 30 
        }

        var monthlyReport []string
        for appliance, energy := range monthlyEnergy {
            monthlyReport = append(monthlyReport, fmt.Sprintf("%s: %.2f kWh", appliance, energy))
        }
        answer = "Konsumsi Energi Per Bulan: " + strings.Join(monthlyReport, ", ")
    }

	if strings.Contains(questionLower, "konsumsi energi tertinggi") || strings.Contains(questionLower, "konsumsi energi terendah") {
		isHighest := strings.Contains(questionLower, "tertinggi")
		targetEnergy := 0.0
		if !isHighest {
			targetEnergy = math.MaxFloat64 // For lowest energy, start with max float value
		}
		targetAppliance := ""

		for i := 0; i < len(table["Appliance"]); i++ {
			energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
			if err != nil {
				log.Printf("Error parsing energy consumption for %s: %v", table["Appliance"][i], err)
				continue
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


	

	if strings.Contains(questionLower, "perangkat di") || strings.Contains(questionLower, "perangkat yang ada di") || strings.Contains(questionLower, "perangkat apa saja yang ada di") || strings.Contains(questionLower, "total perangkat yang berada di") {
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





	if strings.Contains(questionLower, "perangkat dengan status on") || strings.Contains(questionLower, "perangkat dengan status aktif") || strings.Contains(questionLower, "perangkat yang aktif") || strings.Contains(questionLower, "perangkat apa saja yang aktif") {
		activeAppliances := make(map[string]bool) // To ensure uniqueness
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

	if strings.Contains(questionLower, "perangkat dengan status off") || strings.Contains(questionLower, "perangkat dengan status tidak aktif") || strings.Contains(questionLower, "perangkat yang tidak aktif") || strings.Contains(questionLower, "perangkat apa saja yang tidak aktif") {
		inactiveAppliances := make(map[string]bool)
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


	
	// Menangani pertanyaan tentang total penghematan energi
	// Handle total energy cost
	if strings.Contains(questionLower, "total biaya energi") {
		ratePerKWh := 1500.0 // Rate per kWh (in IDR)
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


    // Menangani pertanyaan tentang total biaya energ

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

	if answer == "" {
        log.Println("No matching structure for question. Providing default response.")
        answer = "Pertanyaan tidak dikenali. "
    }

	
	return model.TapasResponse{
		Answer: answer,
		Coordinates: coordinates,  // Tambahkan jika perlu koordinat
		Cells:       cells,        // Tambahkan jika perlu sel
		Aggregator:  aggregator,    // Tambahkan jika perlu agregator
	}, nil

}





