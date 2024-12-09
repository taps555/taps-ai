package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
// Gunakan nama kunci yang lebih deskriptif untuk sesi.
var store = sessions.NewCookieStore([]byte("super-secret-key"))	

// Fungsi untuk mendapatkan sesi
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

func processTableData(table map[string][]string) (model.TapasResponse, error) {
	// Misalnya kita akan menghitung total energi berdasarkan "Appliance"
	applianceEnergy := make(map[string]float64)
	for i := 0; i < len(table["Appliance"]); i++ {
		appliance := table["Appliance"][i]
		energyConsumption, err := strconv.ParseFloat(table["Energy_Consumption"][i], 64)
		if err != nil {
			return model.TapasResponse{}, fmt.Errorf("error parsing energy consumption: %v", err)
		}

		// Menjumlahkan energi per appliance
		applianceEnergy[appliance] += energyConsumption
	}

	// Menyiapkan hasil yang akan dikembalikan
	answer := "Total konsumsi energi perangkat:"
	var coordinates [][]int
	var cells []string
	var aggregator = "SUM"
	
	// Mengonversi data ke format yang diinginkan untuk respons
	for appliance, totalEnergy := range applianceEnergy {
		answer += fmt.Sprintf(" %s: %.2f kWh.", appliance, totalEnergy)
	}

	return model.TapasResponse{
		Answer:      answer,
		Coordinates: coordinates, // Bisa ditambahkan logika jika perlu koordinat spesifik
		Cells:       cells,       // Daftar sel jika diperlukan
		Aggregator:  aggregator,
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
		// Ambil sesi pengguna
		session := getSession(r)
		log.Println("Table from session:", session)
	
		// Tentukan pengaturan sesi
		store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   false, // Nonaktifkan Secure untuk pengujian lokal
		}
	
		// Parse file yang diunggah
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Gagal membaca file yang diunggah", http.StatusBadRequest)
			return
		}
		defer file.Close()
	
		// Membaca konten file
		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Gagal membaca isi file", http.StatusInternalServerError)
			return
		}
	
		// log.Println("File content:", fileContent)
	
		// Proses file menjadi tabel
		table, err := fileService.ProcessFile(string(fileContent))
		if err != nil {
			http.Error(w, "Gagal memproses file: "+err.Error(), http.StatusBadRequest)
			return
		}
	
		// Simpan tabel dalam sesi
		if err := saveSessionData(w, r, table); err != nil {
			log.Printf("Error saving session data: %v\n", err)
			http.Error(w, "Gagal menyimpan data ke sesi", http.StatusInternalServerError)
			return
		}
	
		// Ambil data tabel dari sesi menggunakan getSessionData
		sessionTable, err := getSessionData(r)
		if err != nil {
			log.Printf("Error retrieving session data: %v\n", err)
			http.Error(w, "Gagal mengambil data sesi", http.StatusInternalServerError)
			return
		}
	
		// Proses data tabel dan siapkan format respons
		tapasResponse, err := processTableData(sessionTable)
		if err != nil {
			log.Printf("Error processing table data: %v\n", err)
			http.Error(w, "Gagal memproses data", http.StatusInternalServerError)
			return
		}
	
		// Kirim respons dengan data yang diformat
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "sukses",
			"message": "File berhasil diproses",
			"data": map[string]interface{}{
				"answer": tapasResponse.Answer,
				"coordinates": tapasResponse.Coordinates,
				"cells": tapasResponse.Cells,
				"aggregator": tapasResponse.Aggregator,
			},
		})
	}).Methods("POST")
	









	//2
	// router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
	// 	// Ambil sesi pengguna
	// 	session := getSession(r)
	// 	log.Println("Table from session:", session)

	// 	// Tentukan pengaturan sesi
	// 	store.Options = &sessions.Options{
	// 		Path:     "/",
	// 		MaxAge:   3600,
	// 		HttpOnly: true,
	// 		SameSite: http.SameSiteNoneMode,
	// 		Secure:   false, // Nonaktifkan Secure untuk pengujian lokal
	// 	}

	// 	// Parse file yang diunggah
	// 	file, _, err := r.FormFile("file")
	// 	if err != nil {
	// 		http.Error(w, "Gagal membaca file yang diunggah", http.StatusBadRequest)
	// 		return
	// 	}
	// 	defer file.Close()

	// 	// Membaca konten file
	// 	fileContent, err := io.ReadAll(file)
	// 	if err != nil {
	// 		http.Error(w, "Gagal membaca isi file", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	log.Println("File content:", fileContent)

	// 	// Proses file menjadi tabel
	// 	// Misalnya fileService.ProcessFile mengembalikan map[string][]string
	// 	table := make(map[string][]string)
	// 	// Implementasi pengolahan file sesuai dengan logika yang Anda inginkan
	// 	// table, err := fileService.ProcessFile(string(fileContent)) 

	// 	// Simpan tabel dalam sesi
	// 	if err := saveSessionData(w, r, table); err != nil {
	// 		log.Printf("Error saving session data: %v\n", err)
	// 		http.Error(w, "Gagal menyimpan data ke sesi", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Ambil data tabel dari sesi menggunakan getSessionData
	// 	sessionTable, err := getSessionData(r)
	// 	if err != nil {
	// 		log.Printf("Error retrieving session data: %v\n", err)
	// 		http.Error(w, "Gagal mengambil data sesi", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Proses data tabel dan siapkan format respons
	// 	tapasResponse, err := processTableData(sessionTable)
	// 	if err != nil {
	// 		log.Printf("Error processing table data: %v\n", err)
	// 		http.Error(w, "Gagal memproses data", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Kirim respons dengan data yang diformat
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"status":  "sukses",
	// 		"message": "File berhasil diproses",
	// 		"data": map[string]interface{}{
	// 			"answer": tapasResponse.Answer,
	// 			"coordinates": tapasResponse.Coordinates,
	// 			"cells": tapasResponse.Cells,
	// 			"aggregator": tapasResponse.Aggregator,
	// 		},
	// 	})
	// }).Methods("POST")
	
	
	
	
	
	
	
	

	// File upload endpoint
	// router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
	// 	// Ambil sesi pengguna
	// 	session := getSession(r)
	
	// 	// Tentukan pengaturan sesi (letakkan sebelum penyimpanan sesi)
	// 	store.Options = &sessions.Options{
	// 		Path:     "/",
	// 		MaxAge:   3600,
	// 		HttpOnly: true,
	// 		SameSite: http.SameSiteNoneMode,
	// 		Secure:   false, // Nonaktifkan Secure untuk pengujian lokal
	// 	}
	
	// 	// Parse file yang diunggah
	// 	file, _, err := r.FormFile("file")
	// 	if err != nil {
	// 		http.Error(w, "Gagal membaca file yang diunggah", http.StatusBadRequest)
	// 		return
	// 	}
	// 	defer file.Close()
	
	// 	// Membaca konten file
	// 	fileContent, err := io.ReadAll(file)
	// 	if err != nil {
	// 		http.Error(w, "Gagal membaca isi file", http.StatusInternalServerError)
	// 		return
	// 	}
	
	// 	// Proses file menjadi tabel
	// 	table, err := fileService.ProcessFile(string(fileContent))
	// 	if err != nil {
	// 		http.Error(w, "Gagal memproses file: "+err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	
	// 	// Simpan tabel dalam sesi
	// 	if err := session.Save(r, w); err != nil {
	// 		log.Printf("Error saving session: %v\n", err)  // Log error lebih rinci
	// 		http.Error(w, "Gagal menyimpan sesi", http.StatusInternalServerError)
	// 		return
	// 	}
	
	// 	// Ambil query dari parameter URL (opsional)
	// 	query := r.URL.Query().Get("query")
	// 	if query == "" {
	// 		http.Error(w, "Query diperlukan untuk analisis", http.StatusBadRequest)
	// 		return
	// 	}
	
	// 	// Analisis tabel menggunakan TAPAS
	// 	token := os.Getenv("HUGGINGFACE_TOKEN")

	// 	answer, err := aiService.AnalyzeData(table, query, token)
	// 	if err != nil {
	// 		http.Error(w, "Gagal menganalisis tabel: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	// log.Println("ini resuklt",answer)
	
	// 	// Kirim jawaban ke klien
	// 	// Format ulang jawaban agar lebih menarik
	// 		// Gunakan 'answer' langsung, karena 'result' tidak tersedia
	// 			formattedAnswer := fmt.Sprintf("Hasil analisis: %s", answer)

	// 			// Gunakan formattedAnswer dalam respons
	// 			w.Header().Set("Content-Type", "application/json")
	// 			json.NewEncoder(w).Encode(map[string]interface{}{
	// 				"status":  "success",
	// 				"message": "File berhasil diproses",
	// 				"data": map[string]interface{}{
	// 					"answer": formattedAnswer, // Jawaban yang diformat
	// 				},
	// 			})


	
	// 	// log.Printf("Response headers: %v\n", w.Header())
	// }).Methods("POST")
	
	

	
	
	
	




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









	// router.HandleFunc("/check-session", func(w http.ResponseWriter, r *http.Request) {
	// 	session := getSession(r)
	// 	// Cek apakah data 'table' ada dalam sesi
	// 	table, ok := session.Values["table"]
	// 	if !ok {
	// 		http.Error(w, "Data tabel tidak ditemukan dalam sesi", http.StatusBadRequest)
	// 		return
	// 	}
	
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"table": table,
	// 	})
	// }).Methods("GET")
	

	// Chat endpoint
	// router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
	// 	// Retrieve context and query from the request body
	// 	var request struct {
	// 		Context string `json:"context"`
	// 		Query   string `json:"query"`
	// 	}
	// 	err := json.NewDecoder(r.Body).Decode(&request)
	// 	if err != nil {
	// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 		return
	// 	}

	// 	// Retrieve the session and get the table data
	// 	session := getSession(r)
	// 	table, ok := session.Values["table"].(map[string][]string)
	// 	if !ok {
	// 		http.Error(w, "No table data found in session", http.StatusBadRequest)
	// 		return
	// 	}

	// 	// Analyze the data using AI service
	// 	answer, err := aiService.AnalyzeData(table, request.Query, token)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Return the response
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte(`{"generated_text": "` + answer + `"}`))
	// }).Methods("POST")

	// Enable CORS
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
