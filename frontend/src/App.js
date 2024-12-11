import React from "react";

const UnifiedLayout = () => {
  const styles = {
    container: {
      width: "100vw",
      height: "100vh",
      display: "flex",
      backgroundColor: "#3BA881",
    },
    sidebar: {
      width: "80px",
      height: "98%",
      backgroundColor: "white",
      borderRadius: "35px",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      margin: "10px",
      padding: "20px 10px",
      boxSizing: "border-box",
    },
    icon: {
      width: "40px",
      height: "40px",
      borderRadius: "50%",
      backgroundColor: "#3BA881",
      marginBottom: "20px",
    },
    menuItem: {
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      marginBottom: "20px",
    },
    menuText: {
      color: "#59C19B",
      fontSize: "16px",
      fontFamily: "Outfit, sans-serif",
      fontWeight: "400",
      marginBottom: "10px",
    },
    menuIcon: {
      width: "50px",
      height: "50px",
      backgroundColor: "#59C19B",
      borderRadius: "5px",
    },
    dashboard: {
      flex: 1,
      backgroundColor: "#3BA881",
      padding: "20px",
      boxSizing: "border-box",
      display: "flex",
      flexDirection: "column",
      gap: "20px",
    },
    header: {
      width: "100%",
      display: "grid",
      gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))",
      gap: "20px",
      marginBottom: "20px",
    },
    card: {
      backgroundColor: "white",
      borderRadius: "10px",
      padding: "15px",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
    },
    cardTitle: {
      fontSize: "14px",
      fontWeight: "500",
      color: "#3BA881",
      marginBottom: "10px",
    },
    cardValue: {
      fontSize: "24px",
      fontWeight: "700",
    },
    content: {
      width: "100%",
      backgroundColor: "white",
      borderRadius: "10px",
      padding: "20px",
      display: "grid",
      gridTemplateRows: "1fr 1fr",
      gap: "20px",
      boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
    },
    chart: {
      backgroundColor: "white",
      borderRadius: "10px",
      padding: "15px",
      display: "flex",
      flexDirection: "column",
      justifyContent: "center",
      alignItems: "center",
      boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
    },
  };

  return (
    <div style={styles.container}>
      {/* Sidebar */}
      <div style={styles.sidebar}>
        <div style={styles.icon}></div>
        <div style={styles.menuItem}>
          <span style={styles.menuText}>menu</span>
          <div style={styles.menuIcon}></div>
        </div>
      </div>

      {/* Dashboard */}
      <div style={styles.dashboard}>
        {/* Header Section */}
        <div style={styles.header}>
          <div style={styles.card}>
            <div style={styles.cardTitle}>Result Refrigerator</div>
            <div style={styles.cardValue}>70%</div>
          </div>
          <div style={styles.card}>
            <div style={styles.cardTitle}>Result TV</div>
            <div style={styles.cardValue}>70%</div>
          </div>
          <div style={styles.card}>
            <div style={styles.cardTitle}>Result EVCar</div>
            <div style={styles.cardValue}>40%</div>
          </div>
          <div style={styles.card}>
            <div style={styles.cardTitle}>Result Computer / Laptop</div>
            <div style={styles.cardValue}>150 metric tons CO2/year</div>
          </div>
          <div style={styles.card}>
            <div style={styles.cardTitle}>Result Lighting</div>
            <div style={styles.cardValue}>150 metric tons CO2/year</div>
          </div>
        </div>

        {/* Content Section */}
        <div style={styles.content}>
          <div style={styles.chart}>
            <h3>Energy Usage</h3>
            {/* Replace with your chart component */}
            <p>Chart or graph placeholder</p>
          </div>
          <div style={styles.chart}>
            <h3>Carbon Footprint CO2</h3>
            {/* Replace with your chart component */}
            <p>Chart or graph placeholder</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UnifiedLayout;
