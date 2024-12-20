import React, { useState } from "react";
import { Link } from "react-router-dom";
import ecoElectricIcon from "../assets/images/icon-co.png";
import brain from "../assets/images/brain-circuit.png";
import classNames from "classnames";

const Sidebar = () => {
  const [isCollapsed, setIsCollapsed] = useState(true); // State untuk mengecilkan sidebar
  // Fungsi toggle
  const toggleSidebar = () => {
    setIsCollapsed(!isCollapsed);
  };

  return (
    <div
      className={`h-full bg-white rounded-2xl shadow-xl flex flex-col items-center transition-all duration-300 ${
        isCollapsed ? "w-20 p-2" : "w-60 p-6"
      }`}
      style={{
        background: "linear-gradient(to bottom, #f8f9fa, #e9ecef)",
      }}
    >
      {/* Toggle Button with Hamburger Animation */}
      <div
        className={classNames(
          "cursor-pointer tham tham-e-arrow tham-w-8 tham-h-8 mt-8", // Tambahkan mt-8 di sini
          {
            "tham-active": !isCollapsed, // Aktifkan animasi burger jika collapsed
          }
        )}
        onClick={toggleSidebar}
      >
        <div className="tham-box">
          <div className="tham-inner" />
        </div>
      </div>

      {/* Sidebar Header */}
      <div
        className={`w-full flex items-center justify-center text-white font-bold ${
          isCollapsed ? "hidden" : "block"
        }`}
      ></div>

      {/* Menu Section */}
      <div className="mt-10 w-full">
        {" "}
        {/* Tambahkan mt-10 untuk memberikan jarak */}
        <span
          className={`text-md text-blue-700 font-semibold block mb-3 ${
            isCollapsed ? "hidden" : "block"
          }`}
        >
          Main Menu
        </span>
        <div className="w-full flex items-center mb-4 hover:bg-blue-100 rounded-md p-2 cursor-pointer transition">
          <Link to="/" className="flex items-center w-full">
            <div className="w-10 h-10 bg-blue-500 rounded-md flex items-center justify-center text-white font-bold">
              <img
                src={ecoElectricIcon}
                alt="Energy Insights"
                className="w-11 h-11"
                style={{ objectFit: "contain" }}
              />
            </div>
            {!isCollapsed && (
              <span className="ml-3 text-blue-800 font-medium text-sm">
                Energy Mode
              </span>
            )}
          </Link>
        </div>
        <div className="w-full flex items-center mb-4 hover:bg-blue-100 rounded-md p-2 cursor-pointer transition">
          <Link to="/chat-with-ai" className="flex items-center w-full">
            <div className="w-10 h-10 bg-blue-500 rounded-md flex items-center justify-center text-white font-bold">
              <img
                src={brain}
                alt="Energy Insights"
                className="w-11 h-11"
                style={{ objectFit: "contain" }}
              />
            </div>
            {!isCollapsed && (
              <span className="ml-3 text-blue-800 font-medium text-sm">
                Chat with AI
              </span>
            )}
          </Link>
        </div>
      </div>

      {/* Footer Section */}
      <div className={`mt-auto w-full ${isCollapsed ? "hidden" : "block"}`}>
        <span className="text-xs text-blue-500 font-semibold block mb-2">
          © 2024 Tailored Ai and Power System
        </span>
      </div>
    </div>
  );
};

export default Sidebar;
