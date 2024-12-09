package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
    // Validasi input tabel
    if len(table) == 0 {
        return "", errors.New("tabel kosong")
    }

    // Membuat payload untuk API
    payload := map[string]interface{}{
        "inputs": map[string]interface{}{
            "table": table,
            "query": query,
        },
    }

    // Marshal payload menjadi JSON
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", errors.New("gagal membuat payload JSON")
    }

    // Membuat permintaan HTTP POST
    req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", errors.New("gagal membuat permintaan HTTP")
    }

    // Menambahkan header ke permintaan
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    // Eksekusi permintaan ke API
    resp, err := s.Client.Do(req)
    if err != nil {
        return "", errors.New("gagal mengirim permintaan ke API")
    }
    defer resp.Body.Close()

    // Periksa status kode respons dari API
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("API gagal: %s", string(body))
    }

    // Dekode respons API
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", errors.New("gagal mendekode respons JSON")
    }

    // Debugging: Log respons API
    // log.Printf("API Response: %v", result)

    // Gantilah 'generated_text' dengan 'answer'
    answer, ok := result["answer"].(string)
    if !ok {
        // Jika 'answer' tidak ada, logkan seluruh respons untuk debugging
        log.Printf("Tidak ditemukan 'answer' dalam respons: %v", result)
        return "", errors.New("respon tidak valid atau 'answer' tidak ditemukan")
    }

    // Kembalikan teks yang dihasilkan
    return answer, nil
}


func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	url := "https://api-inference.huggingface.co/models/Qwen/Qwen2.5-Coder-32B-Instruct"

	// Membuat payload request
	payload := map[string]string{
		"inputs": context + "\n" + query,
	}

	// Mengkonversi payload ke JSON
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Membuat request ke API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request
	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Memeriksa status response
	if resp.StatusCode != http.StatusOK {
		return model.ChatResponse{}, fmt.Errorf("failed with status: %s", resp.Status)
	}

	// Membaca response body
	var chatResp []model.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Mengembalikan jawaban
	if len(chatResp) > 0 {
		return chatResp[0], nil
	}
	return model.ChatResponse{}, fmt.Errorf("no response generated")
}