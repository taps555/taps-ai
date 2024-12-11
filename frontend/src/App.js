import React, { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");

  const [response, setResponse] = useState("");
  const [question, setQuestion] = useState("");

  const [responseAI, setResponseAI] = useState("");
  const [questionAI, setQuestionAI] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file) {
      setError("Tidak ada file yang dipilih!");
      return;
    }

    if (!question.trim()) {
      // Memeriksa jika 'question' kosong atau hanya spasi
      setError("Pertanyaan tidak boleh kosong!");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);
    formData.append("question", question); // Menambahkan question ke dalam formData

    try {
      setLoading(true);
      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      // Handle successful response
      console.log("Upload respons:", res.data);
      setResponse(res.data.data); // Set response data from backend
      setLoading(false);
    } catch (error) {
      if (error.response) {
        // Server responded with a status code outside of the 2xx range
        console.error("Error response:", error.response);
        setError(error.response.data || "Terjadi kesalahan pada server.");
      } else if (error.request) {
        // Request was made but no response was received
        console.error("Error request:", error.request);
        setError("Tidak ada respons dari server.");
      } else {
        console.error("Error:", error.message);
        setError("Terjadi kesalahan saat mengirim permintaan.");
      }
    }
  };

  const handleChat = async () => {
    console.log("Sending query:", query); // Debugging output

    try {
      const res = await axios.post("http://localhost:8080/chat", { query });

      if (res.data && res.data.answer) {
        console.log("Respons chat:", res.data.answer);

        // Memisahkan pertanyaan dan jawaban
        const response = res.data.answer;
        const question = query; // Pertanyaan yang dikirimkan
        const answer = response.replace(question, "").trim(); // Menghapus pertanyaan dari jawaban

        // Menyimpan pertanyaan dan jawaban
        setQuestionAI(question); // Menyimpan pertanyaan
        setResponseAI(answer); // Menyimpan jawaban yang sudah dipisahkan

        console.log("Pertanyaan:", question);
        console.log("Jawaban:", answer);
      } else {
        setError("Respons dari server tidak sesuai.");
      }

      setLoading(false); // Adjust based on backend response
    } catch (error) {
      console.error("Error querying chat:", error);
      setError("Terjadi kesalahan saat mengirimkan query.");
      setLoading(false); // Pastikan loading di-set ke false meskipun terjadi error
    }
  };

  if (loading) return <div>Loading...</div>;

  const renderResponseData = () => {
    const { aggregator, answer, cells, coordinates } = response;

    return (
      <div>
        <p>
          <strong>Answer:</strong> {answer}
        </p>
      </div>
    );
  };

  return (
    <div
      style={{
        maxWidth: "600px",
        margin: "0 auto",
        padding: "20px",
        textAlign: "center",
        fontFamily: "Arial, sans-serif",
      }}
    >
      <h1 style={{ color: "#333", marginBottom: "20px" }}>
        Data Analysis Chatbot
      </h1>
      <div style={{ marginBottom: "20px" }}>
        {/* Input file */}
        <input
          type="file"
          onChange={handleFileChange}
          style={{
            padding: "10px",
            marginRight: "10px",
            border: "1px solid #ccc",
            borderRadius: "4px",
          }}
        />
      </div>

      {/* Input untuk pertanyaan */}
      <div style={{ marginBottom: "20px" }}>
        <input
          type="text"
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          placeholder="Masukkan pertanyaan..."
          style={{
            padding: "10px",
            marginBottom: "10px",
            width: "100%",
            border: "1px solid #ccc",
            borderRadius: "4px",
          }}
        />
      </div>

      {/* Tombol upload */}
      <button
        onClick={handleUpload}
        style={{
          padding: "10px 20px",
          backgroundColor: "#007bff",
          color: "white",
          border: "none",
          borderRadius: "4px",
          cursor: "pointer",
        }}
      >
        Upload and Analyze
      </button>

      <div style={{ marginTop: "20px" }}></div>
      <div style={{ marginBottom: "20px" }}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Ask a question..."
          style={{
            padding: "10px",
            marginRight: "10px",
            border: "1px solid #ccc",
            borderRadius: "4px",
            width: "calc(100% - 140px)",
          }}
        />
        <button
          onClick={handleChat}
          style={{
            padding: "10px 20px",
            backgroundColor: "#007bff",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer",
          }}
        >
          Chat
        </button>
      </div>
      <div
        style={{
          marginTop: "20px",
          padding: "10px",
          border: "1px solid #ccc",
          borderRadius: "4px",
          backgroundColor: "#f9f9f9",
        }}
      >
        <h2>Response</h2>
        {response ? renderResponseData() : <p>{error}</p>}
      </div>

      <div
        style={{
          marginTop: "20px",
          padding: "10px",
          border: "1px solid #ccc",
          borderRadius: "4px",
          backgroundColor: "#f9f9f9",
        }}
      >
        <h2>Question AI</h2>
        {questionAI ? <p>{questionAI}</p> : <p>{error}</p>}
        <h2>Response AI</h2>
        {responseAI ? <p>{responseAI}</p> : <p>{error}</p>}
      </div>
    </div>
  );
}

export default App;
