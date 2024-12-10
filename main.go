package main

import (
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var aiService = &service.AIService{Client: &http.Client{}}
// Gunakan nama kunci yang lebih deskriptif untuk sesi.
var store = sessions.NewCookieStore([]byte("super-secret-key"))	
var fileService = &service.FileService{}


func getSession(r *http.Request) *sessions.Session {
	session, err := store.Get(r, "super-secret-key")	
	if err != nil {
		log.Printf("Error getting session: %v\n", err)
	}
	return session
}

func init() {
	// Register the map type before using it in the session
	gob.Register(map[string][]string{})
}

// Fungsi untuk mengambil data sesi
func saveSessionData(w http.ResponseWriter, r *http.Request, table map[string][]string) error {
	session := getSession(r)
	// Simpan tabel dalam sesi
	session.Values["table"] = table
	session.Save(r, w)
	return nil
}

// Fungsi untuk mengambil data dari sesi
func getSessionData(r *http.Request) (map[string][]string, error) {
	session, _ := store.Get(r, "super-secret-key")
	if session.Values["table"] == nil {
		return nil, fmt.Errorf("no table data in session")
	}
	return session.Values["table"].(map[string][]string), nil
}
// Pastikan package strings diimport


// Pastikan package strings diimport



func processTableData(table map[string][]string, question string) (model.TapasResponse, error) {
    var answer string
    applianceEnergy := make(map[string]float64)

    log.Println("Data perangkat yang terdeteksi:")
    for i, appliance := range table["Appliance"] {
        log.Printf("Perangkat %d: %s\n", i+1, appliance)
    }

    questionLower := strings.ToLower(question)

    // Menangani pertanyaan tentang total energi untuk perangkat tertentu
    if strings.Contains(questionLower, "total energi") {
        var applianceQuery string
        // Tentukan perangkat yang disebutkan dalam pertanyaan
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

                if applianceLower == strings.ToLower(applianceQuery) {  // Pastikan perbandingan case-insensitive
                    found = true
                    if _, exists := applianceEnergy[appliance]; !exists {
                        energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
                        if err != nil {
                            return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
                        }
                        applianceEnergy[appliance] += energyConsumption
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

    // Menangani pertanyaan tentang perbandingan energi antara beberapa perangkat
    if strings.Contains(questionLower, "perbandingan energi") {
        var appliances []string
        if strings.Contains(questionLower, "ac") {
            appliances = append(appliances, "AC")
        }
        if strings.Contains(questionLower, "tv") {
            appliances = append(appliances, "TV")
        }
        if strings.Contains(questionLower, "evcar") {
            appliances = append(appliances, "EVCar")
        }
        if strings.Contains(questionLower, "refrigerator") {
            appliances = append(appliances, "Refrigerator")
        } else if strings.Contains(questionLower, "fridge") { // Alias untuk Refrigerator
            appliances = append(appliances, "Refrigerator")
        } else if strings.Contains(questionLower, "televisi") { // Alias untuk TV
            appliances = append(appliances, "TV")
        }

        // Hitung konsumsi energi untuk perangkat yang disebutkan
        comparisonResults := []string{}
        for _, applianceQuery := range appliances {
            found := false
            for i := 0; i < len(table["Appliance"]); i++ {
                appliance := table["Appliance"][i]
                applianceLower := strings.ToLower(appliance)

                if applianceLower == strings.ToLower(applianceQuery) {
                    found = true
                    if _, exists := applianceEnergy[appliance]; !exists {
                        energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
                        if err != nil {
                            return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
                        }
                        applianceEnergy[appliance] += energyConsumption
                    }
                }
            }

            if found {
                comparisonResults = append(comparisonResults, fmt.Sprintf("%s: %.2f kWh", applianceQuery, applianceEnergy[applianceQuery]))
            } else {
                comparisonResults = append(comparisonResults, fmt.Sprintf("%s tidak ditemukan dalam data.", applianceQuery))
            }
        }

        // Gabungkan hasil perbandingan energi
        answer = "Perbandingan Energi Perangkat: " + strings.Join(comparisonResults, ", ")
    }

    // Menangani pertanyaan tentang perangkat yang tidak terdaftar
    if strings.Contains(questionLower, "lampu") || strings.Contains(questionLower, "mesin cuci") || strings.Contains(questionLower, "kipas") {
        answer = "Perangkat yang Anda tanyakan tidak ditemukan dalam data."
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

    // Gabungkan hasil dan kembalikan
    return model.TapasResponse{
        Answer:      answer,
        Coordinates: nil,
        Cells:       nil,
        Aggregator:  "SUM",
    }, nil
}













func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	
	// Set up the router
	router := mux.NewRouter()

	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Batasi ukuran file 10 MB
	
		// Parsing form multipart
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}
	
		// Ambil file yang diunggah
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Gagal membaca file yang diunggah", http.StatusBadRequest)
			return
		}
		defer file.Close()
	
		log.Printf("Uploaded File Name: %s\n", handler.Filename)
		log.Printf("Uploaded File Size: %d\n", handler.Size)
	
		// Baca file CSV menggunakan csv.NewReader
		reader := csv.NewReader(file)
		records, err := reader.ReadAll() // Baca seluruh isi file CSV
		if err != nil {
			log.Printf("Error reading CSV file: %v\n", err)
			http.Error(w, "Invalid CSV file", http.StatusBadRequest)
			return
		}
	
		// Proses data CSV menjadi tabel JSON
		table := make(map[string][]string)
		headers := records[0] // Baris pertama sebagai header
		for _, row := range records[1:] {
			for i, value := range row {
				table[headers[i]] = append(table[headers[i]], value)
			}
		}
	
		// Simpan tabel ke sesi
		if err := saveSessionData(w, r, table); err != nil {
			log.Printf("Error saving session data: %v\n", err)
			http.Error(w, "Gagal menyimpan data ke sesi", http.StatusInternalServerError)
			return
		}
	
		// Ambil pertanyaan dari form
		question := r.FormValue("question")
		if question == "" {
			http.Error(w, "Pertanyaan tidak boleh kosong", http.StatusBadRequest)
			return
		}

		
	
		// Simpan pertanyaan dalam sesi
		session := getSession(r)
		session.Values["question"] = question
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Gagal menyimpan pertanyaan ke sesi", http.StatusInternalServerError)
			return
		}
	
		// Ambil data tabel dari sesi
		sessionTable, err := getSessionData(r)
		if err != nil {
			log.Printf("Error retrieving session data: %v\n", err)
			http.Error(w, "Gagal mengambil data sesi", http.StatusInternalServerError)
			return
		}
	
		// Proses tabel dan pertanyaan
		tapasResponse, err := processTableData(sessionTable, question)
		if err != nil {
			log.Printf("Error processing table data: %v\n", err)
			http.Error(w, "Gagal memproses data", http.StatusInternalServerError)
			return
		}
	
		// Kirim respons dengan data yang diformat
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "sukses",
			"message": "File berhasil diproses dan pertanyaan dijawab.",
			"data": map[string]interface{}{
				"answer":      tapasResponse.Answer,
				"coordinates": tapasResponse.Coordinates,
				"cells":       tapasResponse.Cells,
				"aggregator":  tapasResponse.Aggregator,
			},
		})
	}).Methods("POST")
	
	
	
	







	

	
	
	
	




	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		// Parse the query from the request
		var requestData struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Ensure the user has provided a query
		if requestData.Query == "" {
			http.Error(w, "Query is required", http.StatusBadRequest)
			return
		}

		// Use AI service to get the response
		token := os.Getenv("HUGGINGFACE_TOKEN")
		response, err := aiService.ChatWithAI("", requestData.Query, token)
		if err != nil {
			http.Error(w, "Error processing the query: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the response back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")







	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,  // Mengizinkan cookies
	}).Handler(router)
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
