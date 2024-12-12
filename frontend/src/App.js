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

  useEffect(() => {
    // Set timer to stop loading after 4 seconds
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 1);

    return () => clearTimeout(timer); // Cleanup timer
  }, []);

  if (isLoading) {
    return <LoadingScreen />;
  }

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/chat-with-ai" element={<ChatWithAI />} />
      </Routes>
    </Router>
  );
};

export default App;
