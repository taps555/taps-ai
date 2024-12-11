import React from "react";
import Sidebar from "./sidebar";
import "./style/dashboard.css";

const Dashboard = () => {
  return (
    <div className="dashboard-container">
      {/* Sidebar Component */}
      <Sidebar />

      {/* Main Dashboard Content */}
      <div className="header">
        <div className="card">
          <div className="card-title">Result Refrigerator</div>
          <div className="card-value">70%</div>
        </div>
        <div className="card">
          <div className="card-title">Result TV</div>
          <div className="card-value">70%</div>
        </div>
        <div className="card">
          <div className="card-title">Result EVCar</div>
          <div className="card-value">40%</div>
        </div>
        <div className="card">
          <div className="card-title">Result Computer / Laptop</div>
          <div className="card-value">
            150 <span className="metric">metric tons CO2/year</span>
          </div>
        </div>
        <div className="card">
          <div className="card-title">Result Lighting</div>
          <div className="card-value">
            150 <span className="metric">metric tons CO2/year</span>
          </div>
        </div>

        <div className="content">
          <div className="chart">
            <h3>Energy Usage</h3>
            <p>Chart or graph placeholder</p>
          </div>
          <div className="chart">
            <h3>Carbon Footprint CO2</h3>
            <p>Chart or graph placeholder</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
