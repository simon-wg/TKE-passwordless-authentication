import React, { useState } from "react";
import "./styles.css";
import config from "../config";
import { useNavigate } from "react-router-dom";

const LoginComponent = () => {
  const [username, setUsername] = useState("");
  const [message, setMessage] = useState("");
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (event) => {
    event.preventDefault();
    clearMessages();
    // Send POST Request to client /api/login endpoint
    const response = await fetch(config.clientBaseUrl + "/api/login", {
      method: "POST",
      headers: {
        "Allow-Control-Allow-Origin": "localhost:6060",
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username }),
    });
    if (response.ok) {
      navigate("/loginsuccess");
    } else {
      setError("Failed to sign in user");
    }
  };

  const clearMessages = async () => {
    setMessage("");
    setSuccess("");
    setError("");
  };

  return (
    <div className="container">
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button type="submit">Login</button>
      </form>
      {message && <p className="message">{message}</p>}
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default LoginComponent;
