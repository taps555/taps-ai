import React from "react";
import ReactDOM from "react-dom/client";
import "./components/index.css"; // Import Tailwind CSS here
import App from "./App";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
