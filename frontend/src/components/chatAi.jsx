import React, { useState, useRef, useEffect } from "react";
import Sidebar from "./sidebar";
import axios from "axios";

const ChatWithAI = () => {
  const [query, setQuery] = useState("");
  const [messages, setMessages] = useState([]); // Chat history
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const chatEndRef = useRef(null);

  const cleanResponse = (text) => {
    return text
      .replace(/,/g, "") // Remove commas
      .replace(/\s+/g, " ") // Replace multiple spaces with single space
      .trim(); // Trim leading and trailing spaces
  };

  const handleChat = async () => {
    if (!query.trim()) {
      setError("Query cannot be empty!");
      return;
    }

    setError("");
    setLoading(true);

    // Menambahkan pesan User ke daftar pesan
    setMessages((prev) => [...prev, { sender: "User", text: query }]);

    try {
      const res = await axios.post("http://localhost:8080/chat", { query });

      if (res.data && res.data.answer) {
        const aiResponse = cleanResponse(
          res.data.answer.replace(query, "").trim()
        );

        // Tambahkan placeholder untuk AI response
        setMessages((prev) => [...prev, { sender: "AI", text: "" }]);

        let currentText = "";
        const words = aiResponse.split(" ");
        let index = 0;

        const interval = setInterval(() => {
          if (index < words.length) {
            currentText += (index > 0 ? " " : "") + words[index];
            setMessages((prev) => [
              ...prev.slice(0, prev.length - 1), // Hapus placeholder AI terakhir
              { sender: "AI", text: currentText }, // Tambahkan teks baru
            ]);
            index++;
          } else {
            clearInterval(interval);
          }
        }, 200); // Waktu jeda antara kata (ms)
      } else {
        setError("Unexpected response from the server.");
      }

      setQuery(""); // Kosongkan input
    } catch (err) {
      console.error("Error querying chat:", err);
      setError("Error connecting to AI.");
    } finally {
      setLoading(false);
    }
  };

  // Scroll to the bottom of the chat when new messages are added
  useEffect(() => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300 p-6">
      {/* Sidebar */}
      <Sidebar />

      {/* Main Chat Section */}
      <div className="flex-1 flex flex-col bg-white rounded-2xl  p-8 shadow-xl ml-6  overflow-y-auto animate-fade-in">
        {/* Header */}
        <h1 className="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-800 to-blue-500 mb-6 text-center">
          Chat with AI
        </h1>

        {/* Chat Area */}
        <div className="flex flex-col flex-1 bg-blue-50 rounded-lg p-4 shadow-inner overflow-y-auto">
          {messages.length === 0 && (
            <div className="text-gray-400 italic text-center mt-4">
              Start a conversation by typing your question below.
            </div>
          )}

          {messages.map((msg, index) => (
            <div
              key={index}
              className={`flex ${
                msg.sender === "User" ? "justify-end" : "justify-start"
              } mb-4`}
            >
              <div
                className={`p-4 rounded-lg shadow ${
                  msg.sender === "User"
                    ? "bg-blue-500 text-white"
                    : "bg-gray-100 text-gray-800"
                }`}
                style={{
                  maxWidth: "70%",
                  wordWrap: "break-word",
                  whiteSpace: "pre-wrap", // Preserve line breaks
                }}
              >
                {msg.sender === "AI" && (
                  <p className="text-sm text-gray-700 font-bold mb-1">AI</p> // Darker gray and bold for AI
                )}
                {msg.sender === "User" && (
                  <p className="text-sm text-blue-300 font-extrabold mb-1">
                    You
                  </p> // Bright blue and extra-bold for User
                )}
                <p
                  className={`${
                    msg.sender === "AI"
                      ? "text-gray-800 font-normal"
                      : "text-white font-medium"
                  }`}
                >
                  {msg.text}
                </p>
              </div>
            </div>
          ))}

          <div ref={chatEndRef} />
        </div>

        {/* Input Section */}
        <div className="mt-4 flex">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Type your question here..."
            className="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={handleChat}
            className="px-6 py-2 bg-blue-500 text-white rounded-r-lg hover:bg-blue-600 transition"
            disabled={loading}
            style={{ cursor: loading ? "not-allowed" : "pointer" }}
          >
            {loading ? (
              <svg
                className="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8v8h8a8 8 0 01-8 8v-8H4z"
                ></path>
              </svg>
            ) : (
              "Send"
            )}
          </button>
        </div>

        {/* Error Message */}
        {error && <p className="text-red-500 mt-3">{error}</p>}
      </div>
    </div>
  );
};

export default ChatWithAI;
