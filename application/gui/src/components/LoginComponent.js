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

  const handleGetSignedChallenge = async (event) => {
    event.preventDefault();
    clearMessages();
    setLoading(true);
    // Send POST Request to client /api/login endpoint
    try {
      const response = await fetch(config.clientBaseUrl + "/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json(); // Parse JSON response
      verifySignedChallenge(username, data.signed_challenge);
    } catch (error) {
      setError("Error fetching signed challenge");
      setLoading(false);
    }
  };

  async function verifySignedChallenge(username, signedChallenge) {
    try {
      const response = await fetch("/api/verify", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          username: username,
          signature: signedChallenge,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      navigate("/loginsuccess");
    } catch (error) {
      console.error("Verification failed:", error);
      setError("Error verifying signed challenge!");
      setLoading(false);
    }
  }

  const clearMessages = async () => {
    setMessage("");
    setSuccess("");
    setError("");
  };

  return (
    <div className="container">
      <h2>Login</h2>

      <LoadingCircle loading={loading} />
      <form onSubmit={handleGetSignedChallenge}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button onClick={handleGetSignedChallenge} disabled={loading}>
          {loading ? "Awaiting login" : "Login"}
        </button>
      </form>
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default LoginComponent;
