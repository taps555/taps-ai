import React from "react";
import { Link } from "react-router-dom";
import ecoElectricIcon from "../assets/images/icon-co.png";
import brain from "../assets/images/brain-circuit.png";
import myLogo from "../assets/images/myLogo.png";

const Sidebar = () => {
  return (
    <div
      className="w-64 h-full bg-white rounded-2xl p-6 shadow-xl flex flex-col items-start"
      style={{
        background: "linear-gradient(to bottom, #f8f9fa, #e9ecef)", // Subtle gradient
      }}
    >
      {/* Sidebar Header */}
      <br />

      <div className="w-full rounded-md flex items-center justify-center text-white font-bold">
        <img
          src={myLogo}
          alt="Energy Insights"
          className="w-40 h-70"
          style={{
            objectFit: "contain",
          }}
        />
      </div>

      {/* Menu Section */}
      <br />
      <br />
      <div className="mb-8 w-full">
        <span className="text-md text-blue-700 font-semibold block mb-3 ">
          Main Menu
        </span>
        <div className="w-full flex items-center mb-4 hover:bg-blue-100 rounded-md p-2 cursor-pointer transition">
          <Link to="/" className="flex items-center w-full">
            <div className="w-10 h-10 bg-blue-500 rounded-md flex items-center justify-center text-white font-bold">
              <img
                src={ecoElectricIcon}
                alt="Energy Insights"
                className="w-11 h-11"
                style={{
                  objectFit: "contain",
                }}
              />
            </div>
            <span className="ml-3 text-blue-800 font-medium text-sm">
              Energy Mode
            </span>
          </Link>
        </div>
        <div className="w-full flex items-center mb-4 hover:bg-blue-100 rounded-md p-2 cursor-pointer transition">
          <Link to="/chat-with-ai" className="flex items-center w-full">
            <div className="w-10 h-10 bg-blue-500 rounded-md flex items-center justify-center text-white font-bold">
              <img
                src={brain}
                alt="Energy Insights"
                className="w-11 h-11"
                style={{
                  objectFit: "contain",
                }}
              />
            </div>
            <span className="ml-3 text-blue-800 font-medium text-sm">
              Chat with AI
            </span>
          </Link>
        </div>
      </div>

      {/* Footer Section */}
      <div className="mt-auto w-full">
        <span className="text-xs text-blue-500 font-semibold block mb-2">
          © 2024 Energy Dashboard
        </span>
      </div>
    </div>
  );
};

export default Sidebar;
