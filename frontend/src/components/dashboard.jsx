import React, { useState, useEffect } from "react";
import axios from "axios";
import Sidebar from "./sidebar";

const Dashboard = ({ showPopup }) => {
  const [file, setFile] = useState(null);
  const [response, setResponse] = useState("");
  const [question, setQuestion] = useState("");
  const [applianceData, setApplianceData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [suggestions, setSuggestions] = useState([]);
  const [popupVisible, setPopupVisible] = useState(showPopup);

  const suggestedQuestions = [
    "Total energi semua perangkat",
    "Konsumsi energi per hari",
    "Konsumsi energi per minggu",
    "Perangkat dengan konsumsi energi tertinggi",
    "Perangkat dengan konsumsi energi terendah",
    "Perangkat dengan status on",
    "Perangkat dengan status off",
    "Perangkat di dapur",
    "Perangkat di ruang tamu",
    "Perangkat di kamar tidur",
    "Total penghematan energi",
    "Total biaya energi",
    "lebih tinggi dari 10",
    "lebih rendah dari 10",
  ];

  const energyThresholds = {
    AC: 360,
    Refrigerator: 100,
    "Washing Machine": 100,
    Toaster: 3,
    "LED Lamp": 1.5,
  };

  const removeDuplicates = (input) => {
    const items = input.split(",").map((item) => item.trim()); // Pisahkan berdasarkan koma
    const uniqueItems = [...new Set(items)]; // Menghapus duplikat dengan Set
    return uniqueItems.join(", "); // Gabungkan kembali menjadi string
  };

  const handleInputChange = (value) => {
    setQuestion(value);
    setSuggestions(
      suggestedQuestions.filter((q) =>
        q.toLowerCase().includes(value.toLowerCase())
      )
    );
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      setPopupVisible(false); // Popup otomatis tertutup setelah 3 detik
    }, 3000);

    return () => clearTimeout(timer); // Membersihkan timer saat komponen unmount
  }, []);

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
      setError("");

      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      if (res.data && res.data.data) {
        const { answer, coordinates, cells, aggregator } = res.data.data;

        if (!answer) {
          setError("No relevant data found.");
          setResponse(null);
          setApplianceData([]);
          return;
        }

        const cleanAnswer = removeDuplicates(answer);

        const appliances =
          cleanAnswer.match(/(\w+): ([\d.]+) kWh/g)?.map((item) => {
            const [name, energyString] = item.split(": ");
            const energy = parseFloat(energyString.replace(" kWh", "")) || 0;
            return { name: name?.trim() || "Unknown", energy };
          }) || [];

        const extendedApplianceData = appliances.map((appliance) => ({
          ...appliance,
          dailyEnergy: (appliance.energy / 7).toFixed(2), // Weekly divided by 7
          monthlyEnergy: (appliance.energy * 4).toFixed(2), // Weekly multiplied by 4
        }));

        // Set state dengan jawaban yang sudah bersih
        setResponse({ answer: cleanAnswer, coordinates, cells, aggregator });
        setApplianceData(extendedApplianceData);
      } else {
        setError("Unexpected response format from server.");
      }
      setLoading(false);
    } catch (error) {
      console.error("Error during upload:", error);
      setError("Error uploading file or fetching response.");
      setLoading(false);
    }
  };

  const isAboveThreshold =
    applianceData.energy > (energyThresholds[applianceData.name] || 0);

  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300 p-6 ">
      {/* Sidebar Component */}
      <Sidebar />

      {/* Main Dashboard Content */}

      <div
        className="flex-1 flex flex-col bg-white rounded-2xl p-8 shadow-xl ml-6  animate-fade-in  overflow-y-auto"
        style={{
          height: "auto",
          paddingBottom: "50px",
        }}
      >
        <h2 className="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-800 to-blue-500 mb-6 text-center relative">
          Application Energy Analys
          <span className="absolute -bottom-1 left-1/2 transform -translate-x-1/2 w-20 h-1 bg-gradient-to-r from-blue-800 to-blue-500 rounded-lg"></span>
        </h2>

        {/* Appliance Section */}
        <div className="grid grid-cols-4 gap-6 mb-6">
          {applianceData.map((appliance, index) => {
            const totalEnergy = applianceData.reduce(
              (sum, item) => sum + item.energy,
              0
            );
            const percentage = (appliance.energy / totalEnergy) * 100;
            const dailyEnergy = appliance.energy / 7;
            const monthlyEnergy = appliance.energy * 4;

            // Periksa apakah melebihi batas persentase
            const isAboveThreshold = percentage > 30;

            return (
              <div
                key={index}
                className={`relative rounded-lg p-6 flex flex-col items-center shadow-md hover:shadow-lg transform transition-all duration-500 animate-fade-in-left ${
                  isAboveThreshold
                    ? "bg-red-100 border-2 border-red-500" // Red color if above threshold
                    : "bg-blue-100" // Default color
                }`}
                style={{
                  animationDelay: `${index * 0.1}s`, // Add delay to stagger the animation
                  animationFillMode: "both",
                }}
              >
                {/* Notifikasi Jika Melebihi Threshold */}
                {isAboveThreshold && (
                  <div className="absolute -top-3 -right-3 bg-red-500 text-white text-xs font-bold px-2 py-1 rounded-full shadow-md animate-bounce">
                    ⚠️ High Usage!
                  </div>
                )}

                {/* Nama Perangkat */}
                <h4
                  className={`text-sm font-semibold mb-3 ${
                    isAboveThreshold ? "text-red-600" : "text-blue-600"
                  }`}
                >
                  {appliance.name || "Unknown"}
                </h4>

                <div className="relative w-24 h-24 mb-3">
                  <svg
                    className="absolute top-0 left-0 w-full h-full"
                    viewBox="0 0 36 36"
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
                      className={`${
                        isAboveThreshold ? "text-red-500" : "text-blue-500"
                      }`}
                      strokeWidth="4"
                      strokeDasharray={`${percentage}, 100`}
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
                    <span
                      className={`text-lg font-bold ${
                        isAboveThreshold ? "text-red-600" : "text-blue-800"
                      }`}
                    >
                      {percentage.toFixed(1)}%
                    </span>
                  </div>
                </div>
                <p
                  className={`text-sm font-normal ${
                    isAboveThreshold ? "text-red-600" : "text-blue-800"
                  }`}
                >
                  Total: {appliance.energy.toFixed(2)} kWh
                </p>
                {response?.question?.toLowerCase().includes("daily") && (
                  <p className="text-blue-800 text-sm font-normal">
                    Daily: {dailyEnergy.toFixed(2)} kWh
                  </p>
                )}
                {response?.question?.toLowerCase().includes("monthly") && (
                  <p className="text-blue-800 text-sm font-normal">
                    Monthly: {monthlyEnergy.toFixed(2)} kWh
                  </p>
                )}
              </div>
            );
          })}
        </div>
        {/* File Upload Section */}
        <div className="bg-gray-100 rounded-lg p-6 shadow-lg mb-6 w-full max-w-2xl mx-auto animate-fade-in-up">
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
            onChange={(e) => setFile(e.target.files[0])}
          />
          <div className="relative mt-4">
            <input
              type="text"
              value={question}
              onChange={(e) => handleInputChange(e.target.value)}
              placeholder="Ask about energy usage"
              className="block w-full mt-4 p-3 rounded-md border border-gray-300 text-sm"
            />
            {suggestions.length > 0 && (
              <div className="absolute top-full mt-1 w-full bg-white border border-gray-300 rounded-lg shadow-lg z-10">
                {suggestions.map((suggestion, index) => (
                  <div
                    key={index}
                    onClick={() => {
                      setQuestion(suggestion);
                      setSuggestions([]);
                    }}
                    className="p-2 cursor-pointer hover:bg-blue-100"
                  >
                    {suggestion}
                  </div>
                ))}
              </div>
            )}
          </div>
          <button
            onClick={handleUpload}
            className="mt-4 bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600 transition w-full"
            disabled={loading}
          >
            {loading ? "Processing..." : "Upload and Analyze"}
          </button>
          {error && (
            <div
              className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mt-4"
              role="alert"
            >
              <strong className="font-bold">Error: </strong>
              <span className="block sm:inline">{error}</span>
            </div>
          )}

          {/* Display textual answer */}
          {response?.answer && (
            <div className="flex flex-col items-start mt-6">
              {/* Question Bubble */}
              {question && (
                <div className="bg-gray-200 text-gray-800 p-4 rounded-lg shadow-md mb-4 self-end max-w-2xl">
                  <h4 className="text-sm font-semibold mb-1 text-right">
                    You:
                  </h4>
                  <p className="text-sm">{question}</p>
                </div>
              )}

              {/* Answer or Error Bubble */}
              {response?.answer ? (
                <div className="bg-blue-500 text-white p-4 rounded-lg shadow-md max-w-2xl">
                  <h4 className="text-sm font-semibold mb-1">AI Response:</h4>
                  <p className="text-sm">{response.answer}</p>
                </div>
              ) : (
                <div className="bg-red-500 text-white p-4 rounded-lg shadow-md max-w-2xl">
                  <h4 className="text-sm font-semibold mb-1">AI Response:</h4>
                  <p className="text-sm">
                    Tidak ada pertanyaan yang relevan ditemukan.
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
