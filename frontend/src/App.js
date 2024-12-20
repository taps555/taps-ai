import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./components/dashboard";
import ChatWithAI from "./components/chatAi";

const LoadingScreen = () => (
  <div className="flex items-center justify-center h-screen bg-gradient-to-b from-blue-800 to-blue-500">
    <div className="text-white text-center space-y-4">
      <h1 className="text-3xl font-bold animate-bounce">
        Hallo, Selamat Datang di <span className="text-yellow-400">TAPS</span>
      </h1>
      <p className="text-lg">Mengoptimalkan Energi untuk Masa Depan</p>
      <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-yellow-400 mx-auto"></div>
    </div>
  </div>
);

const App = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [showGuide, setShowGuide] = useState(false);
  const [currentSlide, setCurrentSlide] = useState(0);
  const [showDashboardPopup, setShowDashboardPopup] = useState(false);

  const slides = [
    {
      title: "Upload CSV Anda!",
      description:
        "Mulai dengan mengunggah file CSV Anda. Pastikan formatnya sesuai, dan kami akan membantu menganalisis energi perangkat Anda secara efisien.",
      image: "/images/csv.png", // Ganti dengan path gambar Anda
    },
    {
      title: "Input Pertanyaan Anda",
      description:
        "Gunakan kolom input untuk menanyakan apa saja terkait konsumsi energi. Anda juga dapat memilih pertanyaan yang disarankan.",
      image: "/images/image.png", // Ganti dengan path gambar Anda
    },
    {
      title: "Gunakan Chat AI",
      description:
        "Berinteraksi dengan AI untuk mendapatkan analisis canggih dan jawaban atas pertanyaan Anda tentang konsumsi energi.",
      image: "/images/Screenshot 2024-12-18 172508.png", // Ganti dengan path gambar Anda
    },
  ];

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false);
      setShowGuide(true);
    }, 4000);

    return () => clearTimeout(timer);
  }, []);

  const nextSlide = () => {
    if (currentSlide < slides.length - 1) {
      setCurrentSlide(currentSlide + 1);
    } else {
      setShowGuide(false); // Tutup panduan jika slide terakhir
      setShowDashboardPopup(true);
    }
  };

  const previousSlide = () => {
    if (currentSlide > 0) {
      setCurrentSlide(currentSlide - 1);
    }
  };

  if (isLoading) {
    return <LoadingScreen />;
  }

  return (
    <Router>
      {showGuide && (
        <div className="fixed inset-0 backdrop-blur-md bg-black bg-opacity-40 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-8 w-11/12 max-w-xl shadow-lg relative animate-fade-in">
            {/* Konten Slide */}
            <div className="text-center">
              <h2 className="text-2xl font-bold mb-4 text-blue-800">
                {slides[currentSlide].title}
              </h2>
              <p className="text-gray-700 mb-4">
                {slides[currentSlide].description}
              </p>
              <img
                src={slides[currentSlide].image}
                alt={slides[currentSlide].title}
                className="w-full h-48 object-contain mb-6"
              />
            </div>

            {/* Tombol Navigasi */}
            <div className="flex justify-between">
              {currentSlide > 0 && (
                <button
                  onClick={previousSlide}
                  className="bg-gray-300 text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-400 transition"
                >
                  Kembali
                </button>
              )}
              <button
                onClick={nextSlide}
                className={`ml-auto px-4 py-2 rounded-lg text-white ${
                  currentSlide === slides.length - 1
                    ? "bg-green-500 hover:bg-green-600"
                    : "bg-blue-500 hover:bg-blue-600"
                } transition`}
              >
                {currentSlide === slides.length - 1
                  ? "Mulai Aplikasi"
                  : "Lanjut"}
              </button>
            </div>
          </div>
        </div>
      )}

      <Routes>
        <Route
          path="/"
          element={<Dashboard showPopup={showDashboardPopup} />}
        />

        <Route path="/chat-with-ai" element={<ChatWithAI />} />
      </Routes>
    </Router>
  );
};

export default App;
