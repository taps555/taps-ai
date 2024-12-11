import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./components/dashboard";

const App = () => {
  return (
    <Router>
      {/* Container for Full-Page Layout */}

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
    </Router>
  );
};

export default App;
