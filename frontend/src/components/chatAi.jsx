import React, { useState } from "react";
import axios from "axios";
import Sidebar from "./sidebar"; // Assuming Sidebar is already created and styled

const ChatWithAI = () => {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");
  const [question, setQuestion] = useState("");
  const [applianceData, setApplianceData] = useState([]);
  const [responseAI, setResponseAI] = useState("");
  const [questionAI, setQuestionAI] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
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

  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300">
      {/* Sidebar */}
      <Sidebar />

      {/* Main Chat Section */}
      <div className="flex-1 flex flex-col bg-white rounded-2xl p-6 shadow-xl mx-6 my-6">
        {/* Header */}
        <h1 className="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-800 to-blue-500 mb-6 text-center">
          Chat with AI
        </h1>

        {/* Chat Area */}
        <div className="flex flex-col flex-1 bg-blue-50 rounded-lg p-4 shadow-inner overflow-hidden">
          {/* Chat Messages */}
          <div className="flex-1 overflow-y-auto px-4 py-2">
            {response ? (
              <div className="bg-white p-4 rounded-lg shadow mb-4">
                <p className="text-gray-800">
                  <strong>AI:</strong> {response}
                </p>
              </div>
            ) : (
              <div className="text-gray-400 italic text-center">
                Your AI response will appear here...
              </div>
            )}
          </div>

          {/* Input Section */}
          <div className="flex mt-4">
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
        </div>

        {/* Error Message */}
        {error && <p className="text-red-500 mt-3">{error}</p>}
      </div>
    </div>
  );
};

export default ChatWithAI;
