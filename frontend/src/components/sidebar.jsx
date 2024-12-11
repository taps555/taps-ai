import React from "react";
import { Link } from "react-router-dom";

const Sidebar = () => {
  return (
    <div className="w-20 h-full bg-white rounded-2xl p-4 shadow-lg flex flex-col items-center">
      <h3 className="text-green-500 font-bold text-sm mb-5">Energy Mode</h3>
      <div className="mb-4">
        <span className="text-xs text-green-700 block mb-2">Menu</span>
        <div className="w-10 h-10 bg-green-300 rounded-md"></div>
      </div>
      <div>
        <span className="text-xs text-green-700 block mb-2">Chat with AI</span>
        <Link to="/chat-with-ai">
          <div className="w-10 h-10 bg-green-300 rounded-md"></div>
        </Link>
      </div>
    </div>
  );
};

export default Sidebar;
