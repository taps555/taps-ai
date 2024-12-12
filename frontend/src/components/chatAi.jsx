import React, { useState } from "react";
import Sidebar from "./sidebar";

const ChatWithAI = () => {
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChat = async () => {
    if (!query.trim()) {
      setError("Query cannot be empty!");
      return;
    }

    setError("");
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query }),
      });
      const data = await res.json();

      if (data.answer) {
        setResponse(data.answer);
      } else {
        setError("No response from AI.");
      }
    } catch (err) {
      setError("Error connecting to AI.");
    }

    setLoading(false);
  };

  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300 p-6">
      {/* Sidebar */}
      <Sidebar />

      {/* Main Content */}
      <div
        className="flex-1 flex flex-col bg-white rounded-2xl p-8 shadow-xl mx-6 my-6"
        style={{
          height: "auto",
          paddingBottom: "50px", // Ensures some padding at the bottom for spacing
        }}
        >
        <h2 className="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-800 to-blue-500 mb-6 text-center relative">
          Chat with AI
          <span className="absolute -bottom-1 left-1/2 transform -translate-x-1/2 w-20 h-1 bg-gradient-to-r from-blue-800 to-blue-500 rounded-lg"></span>
        </h2>

        <div className="flex-1 flex flex-col items-center justify-between bg-blue-100 rounded-lg p-6 shadow-lg">
          {/* Chat Display */}
          <div className="flex-1 w-full overflow-y-auto mb-6 bg-white rounded-lg shadow-inner p-4">
            {response ? (
              <p className="text-gray-800">{response}</p>
            ) : (
              <p className="text-gray-400 italic">
                Your AI response will appear here...
              </p>
            )}
          </div>

          {/* Query Input */}
          <div className="flex items-center w-full">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Type your question here..."
              className="flex-1 px-4 py-2 rounded-l-lg border border-gray-300 text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={handleChat}
              className="px-6 py-2 bg-blue-500 text-white rounded-r-lg hover:bg-blue-600 transition"
              disabled={loading}
            >
              {loading ? "Sending..." : "Send"}
            </button>
          </div>

          {/* Error Message */}
          {error && <p className="text-red-500 mt-3">{error}</p>}
        </div>
      </div>
    </div>
  );
};

export default ChatWithAI;
