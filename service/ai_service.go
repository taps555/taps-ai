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
    if len(table) == 0 {
        return "", errors.New("tabel kosong")
    }

    payload := map[string]interface{}{
        "inputs": map[string]interface{}{
            "table": table,
            "query": query,
        },
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", errors.New("gagal membuat payload JSON")
    }

    req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", errors.New("gagal membuat permintaan HTTP")
    }

    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.Client.Do(req)
    if err != nil {
        return "", errors.New("gagal mengirim permintaan ke API")
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("API gagal: %s", string(body))
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", errors.New("gagal mendekode respons JSON")
    }

    // Cek apakah 'answer' ada dalam respons
    answer, ok := result["answer"].(string)
    if !ok {
        log.Printf("Tidak ditemukan 'answer' dalam respons: %v", result)
        return "", errors.New("respon tidak valid atau 'answer' tidak ditemukan")
    }

    // Memeriksa apakah ada confidence score atau data lain yang bisa digunakan untuk meningkatkan jawaban
    confidence, _ := result["confidence"].(float64)
    if confidence < 0.5 {
        log.Printf("Tingkat kepercayaan rendah: %.2f", confidence)
    }

    return answer, nil
}


func (s *AIService) ChatWithAI(context, query, token string) (string, error) {
	url := "https://api-inference.huggingface.co/models/Qwen/Qwen2.5-Coder-32B-Instruct"

	// Membuat payload request dengan max_new_tokens
	payload := map[string]interface{}{
		"inputs": context + "\n" + query,
		"parameters": map[string]interface{}{
			"max_new_tokens": 550,  // Menentukan jumlah token maksimal untuk respons
			"temperature":    0.7,  // Menentukan suhu untuk variasi respons//
			"streaming":      true, // Mengaktifkan streaming untuk respons bertahap
		},
	}

	// Mengkonversi payload ke JSON
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Membuat request ke API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Memeriksa status response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed with status: %s", resp.Status)
	}

	// Membaca response body
	var chatResp []model.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Memastikan ada respons
	if len(chatResp) > 0 {
		// Mengembalikan hanya jawaban dari model AI
		return chatResp[0].GeneratedText, nil
	}

	return "", fmt.Errorf("no response generated")
}

