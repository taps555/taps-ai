import React, { useState, useRef, useEffect } from "react";
import Sidebar from "./sidebar";

const ChatWithAI = () => {
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [responseAI, setResponseAI] = useState(false);

  const responseRef = useRef(null);

  const handleChat = async () => {
    if (!query.trim()) {
      setError("Query cannot be empty!");
      return;
    }

    setError("");
    setLoading(true);

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

  useEffect(() => {
    if (responseRef.current) {
      responseRef.current.scrollTop = responseRef.current.scrollHeight;
    }
  }, [response]);

  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300 p-6">
      {/* Sidebar */}
      <Sidebar />

      {/* Main Content */}
      <div className="flex-1 flex flex-col bg-white rounded-2xl p-8 shadow-xl mx-6 my-6">
        <h2 className="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-800 to-blue-500 mb-6 text-center">
          Chat with AI
        </h2>

        <div className="flex-1 flex flex-col items-center justify-between bg-blue-100 rounded-lg p-6 shadow-lg">
          <div
            ref={responseRef}
            className="flex-1 w-full overflow-y-auto mb-6 bg-white rounded-lg shadow-inner p-4"
          >
            {response ? (
              <p className="text-gray-800">{response}</p>
            ) : (
              <p className="text-gray-400 italic">
                Your AI response will appear here...
              </p>
            )}
          </div>

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

          {error && <p className="text-red-500 mt-3">{error}</p>}
        </div>
      </div>
    </div>
  );
};

export default ChatWithAI;
