import React, { useState } from "react";
import "./styles.css";
import config from "../config";
import { useNavigate } from "react-router-dom";
import LoadingCircle from "./LoadingCircle";


const LoginComponent = () => {
  const [username, setUsername] = useState("");
  const [message, setMessage] = useState("");
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async (event) => {
    event.preventDefault();
    clearMessages();
    setLoading(true);
    // Send POST Request to localhost:8080/api/initialize-login
    const response = await fetch("http://localhost:8080/api/initialize-login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username }),
    });
    setLoading(false);
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

      <LoadingCircle loading={loading} />
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button onClick={handleLogin} disabled={loading}>
          {loading ? "Awaiting login" : "Login"}
        </button>
      </form>
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
}


export default LoginComponent;
