import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./components/dashboard";
import Sidebar from "./components/sidebar";

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Sidebar />} />
        <Route path="/chat-with-ai" element={<div>Chat with AI Page</div>} />
      </Routes>
    </Router>
  );
};

export default App;
