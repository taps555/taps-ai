import React, { useState } from "react";
import Sidebar from "./sidebar";

const Dashboard = () => {
  const [uploadedFile, setUploadedFile] = useState(null);

  const handleFileUpload = (event) => {
    setUploadedFile(event.target.files[0]);
  };

  return (
    <div className="flex h-screen bg-gradient-to-r from-blue-500 to-blue-300 p-8">
      {/* Sidebar Component */}
      <Sidebar />

      {/* Main Dashboard Content */}
      <div className="flex-1 flex flex-col bg-white rounded-2xl p-8 shadow-xl ml-8">
        {/* Header Section */}
        <div className="grid grid-cols-5 gap-6 mb-6">
          <div className="bg-gray-300  rounded-lg p-5 flex flex-col items-center shadow-md hover:shadow-lg transition">
            <h4 className="text-blue-600 text-sm font-semibold mb-2">
              Refrigerator
            </h4>
            <p className="text-blue-800 text-2xl font-bold">70%</p>
          </div>
          <div className="bg-gray-300  rounded-lg p-5 flex flex-col items-center shadow-md hover:shadow-lg transition">
            <h4 className="text-blue-600 text-sm font-semibold mb-2">TV</h4>
            <p className="text-blue-800 text-2xl font-bold">70%</p>
          </div>
          <div className="bg-gray-300  rounded-lg p-5 flex flex-col items-center shadow-md hover:shadow-lg transition">
            <h4 className="text-blue-600 text-sm font-semibold mb-2">EVCar</h4>
            <p className="text-blue-800 text-2xl font-bold">40%</p>
          </div>
          <div className="bg-gray-300  rounded-lg p-5 flex flex-col items-center shadow-md hover:shadow-lg transition">
            <h4 className="text-blue-600 text-sm font-semibold mb-2">
              Computer
            </h4>
            <p className="text-blue-800 text-2xl font-bold">
              150{" "}
              <span className="text-xs font-normal">metric tons CO2/year</span>
            </p>
          </div>
          <div className="bg-gray-300  rounded-lg p-5 flex flex-col items-center shadow-md hover:shadow-lg transition">
            <h4 className="text-blue-600 text-sm font-semibold mb-2">
              Lighting
            </h4>
            <p className="text-blue-800 text-2xl font-bold">
              150{" "}
              <span className="text-xs font-normal">metric tons CO2/year</span>
            </p>
          </div>
        </div>

        {/* Content Area */}
        <div className="grid grid-cols-2 gap-6 mb-6">
          <div className="bg-blue-50 rounded-lg p-6 shadow-md flex items-center justify-center">
            <h3 className="text-blue-600 text-lg font-semibold">
              Energy Usage
            </h3>
          </div>
          <div className="bg-blue-50 rounded-lg p-6 shadow-md flex items-center justify-center">
            <h3 className="text-blue-600 text-lg font-semibold">
              Carbon Footprint CO2
            </h3>
          </div>
        </div>

        {/* File Upload Section */}
        <div className="bg-blue-50 rounded-lg p-6 shadow-md">
          <h3 className="text-lg font-semibold text-blue-600 mb-3">
            Upload a File
          </h3>
          <input
            type="file"
            className="block w-3/4 text-sm text-blue-600 file:mr-4 file:py-2 file:px-4
                       file:rounded-md file:border-0
                       file:text-sm file:font-semibold
                       file:bg-gray-300  file:text-blue-800
                       hover:file:bg-blue-200 transition"
            onChange={handleFileUpload}
          />
          {uploadedFile && (
            <p className="mt-3 text-sm text-blue-600">
              Uploaded file: <strong>{uploadedFile.name}</strong>
            </p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
