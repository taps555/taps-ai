import React, { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");
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

    const formData = new FormData();
    formData.append("file", file);

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
      console.error("Error uploading file:", error);
      setError("Terjadi kesalahan saat meng-upload file, coba lagi.");
      setLoading(false);
    }
  };

  const handleChat = async () => {
    console.log("Sending query:", query); // Debugging output

    try {
      const res = await axios.post("http://localhost:8080/chat", { query });
      if (res.data && res.data.generated_text) {
        console.log("Respons chat:", res.data);
        setResponse(res.data.generated_text); // Memperbarui respons dengan teks yang dihasilkan
      } else {
        setError("Respons dari server tidak sesuai.");
      }
      setLoading(false); // Adjust based on backend response
    } catch (error) {
      console.error("Error querying chat:", error);
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
        <p>
          <strong>Aggregator:</strong> {aggregator}
        </p>
        <p>
          <strong>Cells:</strong>{" "}
          {cells && cells.length > 0
            ? cells.join(", ")
            : "No cells data available"}
        </p>
        <p>
          <strong>Coordinates:</strong>{" "}
          {coordinates && coordinates.length > 0
            ? coordinates
                .map((coord) => `(${coord[0]}, ${coord[1]})`)
                .join(", ")
            : "No coordinates data available"}
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
      </div>
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
        {response ? (
          renderResponseData()
        ) : (
        )}
      </div>
    </div>
  );
}

export default App;
