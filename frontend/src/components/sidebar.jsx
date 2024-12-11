import React from "react";
import { Link } from "react-router-dom";
import "./style/sidebar.css";
import Dashboard from "./dashboard";

const Sidebar = () => {
  return (
    <div className="sidebar">
      <Dashboard />
      <h3 className="sidebar-title">Energy Mode</h3>
      <div className="menu-item">
        <span className="menu-text">Menu</span>
        <div className="menu-icon"></div>
      </div>
      <div className="menu-item">
        <span className="menu-text">Chat with AI</span>
        <Link to="/chat-with-ai">
          <div className="menu-icon"></div>
        </Link>
      </div>
    </div>
  );
};

export default Sidebar;
