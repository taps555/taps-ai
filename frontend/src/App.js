import React, { useState } from "react";
import axios from "axios";
import "./frontend/src/app.css"; // Import the CSS file

const UnifiedLayout = () => {
  const [file, setFile] = useState(null);
  const [question, setQuestion] = useState("");
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");
  const [responseAI, setResponseAI] = useState("");
  const [questionAI, setQuestionAI] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file || !question.trim()) {
      setError("File and question are required.");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);
    formData.append("question", question);

    try {
      setLoading(true);
      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      setResponse(res.data.data);
      setError("");
    } catch (err) {
      setError("Error uploading file.");
    } finally {
      setLoading(false);
    }
  };

  const handleChat = async () => {
    if (!query.trim()) {
      setError("Query is required.");
      return;
    }

    try {
      setLoading(true);
      const res = await axios.post("http://localhost:8080/chat", { query });
      setQuestionAI(query);
      setResponseAI(res.data.answer);
      setError("");
    } catch (err) {
      setError("Error fetching chat response.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      {/* Sidebar */}
      <div className="sidebar">
        <div className="icon"></div>
        <div className="menu-item">
          <span className="menu-text">menu</span>
          <div className="menu-icon"></div>
        </div>
      </div>

      {/* Dashboard */}
      <div className="dashboard">
        {/* Header Section */}
        <div className="header">
          <div className="card">
            <div className="card-title">Upload File</div>
            <input type="file" onChange={handleFileChange} />
            <input
              type="text"
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder="Enter question..."
            />
            <button onClick={handleUpload}>Upload</button>
          </div>
          <div className="card">
            <div className="card-title">Chat</div>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Ask a question..."
            />
            <button onClick={handleChat}>Send</button>
          </div>
        </div>

        {/* Content Section */}
        <div className="content">
          <div className="chart">
            <h3>Response</h3>
            {response ? (
              <p>{JSON.stringify(response)}</p>
            ) : (
              <p>{error || "No response yet."}</p>
            )}
          </div>
          <div className="chart">
            <h3>AI Chat</h3>
            <p>
              <strong>Question:</strong> {questionAI || "No question yet."}
            </p>
            <p>
              <strong>Response:</strong> {responseAI || "No response yet."}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UnifiedLayout;
