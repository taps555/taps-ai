import React, { useState } from "react";
import axios from "axios";
import Sidebar from "./sidebar";

const Dashboard = () => {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");
  const [question, setQuestion] = useState("");
  const [applianceData, setApplianceData] = useState([]);
  const [responseAI, setResponseAI] = useState("");
  const [questionAI, setQuestionAI] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [uploadedFile, setUploadedFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const calculatePercentage = (value, max = 100) =>
    Math.min((value / max) * 100, 100);

  const handleUpload = async () => {
    if (!file) {
      setError("No file selected!");
      return;
    }

    if (!question.trim()) {
      setError("Question cannot be empty!");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);
    formData.append("question", question);

    try {
      setLoading(true);
      setError(""); // Clear previous errors

      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      if (res.data && res.data.data) {
        const { answer, coordinates, cells, aggregator } = res.data.data;

        // Parse appliance data from the answer (e.g., "Refrigerator: 201.60 kWh, TV: 28.00 kWh")
        const appliances = answer.split(", ").map((item) => {
          const [name, energyString] = item.split(": ");
          const energy = parseFloat(energyString?.replace(" kWh", "")) || 0;
          return { name: name?.trim() || "Unknown", energy };
        });

        // Update states
        setResponse({ answer, coordinates, cells, aggregator });
        setApplianceData(appliances); // Update appliance data for progress bars
        setLoading(false);
      } else {
        setError("Unexpected response format from server.");
        setLoading(false);
      }
    } catch (error) {
      console.error("Error during upload:", error);
      setError("Error uploading file or fetching response.");
      setLoading(false);
    }
  };

  // const handleChat = async () => {
  //   if (!query.trim()) {
  //     setError("Query cannot be empty!");
  //     return;
  //   }

  //   try {
  //     const res = await axios.post("http://localhost:8080/chat", { query });

  //     if (res.data && res.data.answer) {
  //       const response = res.data.answer;
  //       const question = query;
  //       const answer = response.replace(question, "").trim();

  //       setQuestionAI(question);
  //       setResponseAI(answer);
  //     } else {
  //       setError("Unexpected response from server.");
  //     }

  //     setLoading(false);
  //   } catch (error) {
  //     setError("Error querying chat.");
  //     setLoading(false);
  //   }
  // };
  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300 p-6">
      {/* Sidebar Component */}
      <Sidebar />

      {/* Main Dashboard Content */}
      <div className="flex-1 flex flex-col bg-white rounded-2xl p-8 shadow-xl ml-6">
        <h2 className="text-blue-800 text-2xl font-semibold mb-6 text-center">
          Application Energy Analys
        </h2>

        {/* Appliance Section */}
        <div className="grid grid-cols-3 gap-6 mb-6">
          {loading ? (
            <div className="flex justify-center items-center col-span-3">
              <div className="loader">Loading...</div>
            </div>
          ) : (
            applianceData.map((appliance, index) => {
              const energyValue =
                isNaN(appliance.energy) || appliance.energy === null
                  ? 0
                  : appliance.energy;
              const normalizedEnergy = Math.min(100, energyValue);
              const displayEnergy =
                energyValue > 100
                  ? `${energyValue.toFixed(1)} kWh (split)`
                  : `${energyValue.toFixed(1)} kWh`;

              return (
                <div
                  key={index}
                  className="bg-blue-100 rounded-lg p-6 flex flex-col items-center shadow-md hover:shadow-lg transition"
                >
                  <h4 className="text-blue-600 text-sm font-semibold mb-3">
                    {appliance.name || "Unknown"}
                  </h4>

                  {/* Circular Progress Bar */}
                  <div className="relative w-24 h-24 mb-3">
                    <svg
                      className="absolute top-0 left-0 w-full h-full"
                      viewBox="0 0 36 36"
                      aria-label={`${appliance.name}: ${
                        (normalizedEnergy / 100) * 100
                      }%`}
                    >
                      <circle
                        className="text-blue-200"
                        strokeWidth="4"
                        stroke="currentColor"
                        fill="transparent"
                        r="16"
                        cx="18"
                        cy="18"
                      />
                      <circle
                        className="text-blue-500 progress-circle"
                        strokeWidth="4"
                        strokeDasharray={`${
                          (normalizedEnergy / 100) * 100
                        }, 100`}
                        strokeLinecap="round"
                        stroke="currentColor"
                        fill="transparent"
                        r="16"
                        cx="18"
                        cy="18"
                        style={{
                          transform: "rotate(-90deg)",
                          transformOrigin: "50% 50%",
                        }}
                      />
                    </svg>
                    <div className="absolute inset-0 flex items-center justify-center">
                      <span className="text-blue-800 text-lg font-bold">
                        {((normalizedEnergy / 100) * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>

                  <p className="text-blue-800 text-sm font-normal">
                    {displayEnergy}
                  </p>
                </div>
              );
            })
          )}
        </div>

        {/* File Upload Section */}
        <div className="bg-gray-100 rounded-lg p-6 shadow-lg mb-6 w-full max-w-2xl mx-auto">
          <h3 className="text-lg font-semibold text-gray-800 mb-3 text-center">
            Upload a File
          </h3>
          <input
            type="file"
            className="block w-full text-sm text-gray-600 file:mr-4 file:py-2 file:px-4
                       file:rounded-md file:border-0
                       file:text-sm file:font-semibold
                       file:bg-gray-300 file:text-gray-800
                       hover:file:bg-gray-400 transition"
            onChange={handleFileChange}
          />
          <input
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            placeholder="Ask about energy usage (e.g., 'total energi AC')"
            className="block w-full mt-4 p-3 rounded-md border border-gray-300 text-sm"
          />
          <button
            onClick={handleUpload}
            className="mt-4 bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600 transition w-full"
            disabled={loading}
          >
            {loading ? "Processing..." : "Upload and Analyze"}
          </button>
          {error && <p className="text-red-500 mt-3 text-center">{error}</p>}
        </div>

        {/* Detailed Response Section */}
        {response && (
          <div className="bg-gray-100 rounded-lg p-6 shadow-lg w-full max-w-2xl mx-auto">
            <h3 className="text-lg font-semibold text-gray-800 mb-3 text-center">
              Detailed Response
            </h3>
            <p className="text-gray-700">
              <strong>Answer:</strong> {response.answer}
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
