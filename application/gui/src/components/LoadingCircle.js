import React from "react";
import './LoadingCircle.css';

const LoadingCircle = ({ loading }) => {
  return (
    <div className={`image-container ${loading ? "loading" : ""}`}>
    <img
      src="https://tillitis.se/content/uploads/2023/09/tkey-case-1024x391.png"
      alt="TKey Circle"
    />
    <div className="spinner"></div>
  </div>
  );
};

export default LoadingCircle;
