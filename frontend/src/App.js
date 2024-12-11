import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./components/dashboard";

const App = () => {
  return (
    <Router>
      {/* Container for Full-Page Layout */}
      <div className="w-screen h-screen bg-green-500 flex">
        <Routes>
          {/* Main Routes */}
          <Route path="/" element={<Dashboard />} />
          <Route
            path="/chat-with-ai"
            element={
              <div className="flex flex-col items-center justify-center text-white text-lg font-bold">
                Chat with AI Page
              </div>
            }
          />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
