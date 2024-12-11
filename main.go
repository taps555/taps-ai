package main

import (
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var aiService = &service.AIService{Client: &http.Client{}}
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
	gob.Register(map[string][]string{})
}

func saveSessionData(w http.ResponseWriter, r *http.Request, table map[string][]string) error {
	session := getSession(r)
	// Simpan tabel dalam sesi
	session.Values["table"] = table
	session.Save(r, w)
	return nil
}

func getSessionData(r *http.Request) (map[string][]string, error) {
	session, _ := store.Get(r, "super-secret-key")
	if session.Values["table"] == nil {
		return nil, fmt.Errorf("no table data in session")
	}
	return session.Values["table"].(map[string][]string), nil
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}


	router := mux.NewRouter()

	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Batasi ukuran file 10 MB

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
		tapasResponse, err := fileService.ProcessTableData(sessionTable, question)
		if err != nil {
			log.Printf("Error processing table data: %v\n", err)
			http.Error(w, "Gagal memproses data", http.StatusInternalServerError)
			return
		}

		log.Printf("TapasResponse: %+v\n", tapasResponse)
	
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
		var requestData struct {
			Query string `json:"query"`
		}
	
		// Mengambil data dari request body
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	
		// Memastikan query tidak kosong
		if requestData.Query == "" {
			http.Error(w, "Query is required", http.StatusBadRequest)
			return
		}
	
		// Menggunakan layanan AI untuk mendapatkan respons
		token := os.Getenv("HUGGINGFACE_TOKEN")
		answer, err := aiService.ChatWithAI("", requestData.Query, token)
		if err != nil {
			http.Error(w, "Error processing the query: "+err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Menyiapkan data respons dengan query dan answer terpisah
		responseData := struct {
			Query  string `json:"query"`
			Answer string `json:"answer"`
		}{
			Query:  requestData.Query,
			Answer: answer,
		}
	
		// Mengirimkan respons kembali ke klien dalam format JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseData)
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
