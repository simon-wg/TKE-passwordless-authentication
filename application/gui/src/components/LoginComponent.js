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

  /**
   * Handles the process of getting a signed challenge from the server.
   *
   * This function sends a POST request to the /api/login endpoint with the username
   * to get a signed challenge. If the request is successful, it calls the
   * verifySignedChallenge function with the username and the signed challenge.
   *
   * @param {Event} event - The event object from the form submission.
   * @throws {Error} - Throws an error if the HTTP request fails.
   */

  const handleGetSignedChallenge = async (event) => {
    event.preventDefault();
    clearMessages();
    setLoading(true);
    setMessage("Please touch your TKey to proceed with the login.");
    try {
      const response = await fetch(config.clientBaseUrl + "/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username }),
      });

      if (!response.ok) {
        var errorMessage = await response.text();
        setError(errorMessage);
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      verifySignedChallenge(username, data.signed_challenge);
    } catch (error) {
      setLoading(false);
    }
  };

  /**
   * Verifies the signed challenge for the given username.
   *
   * @param {string} username - The username of the user.
   * @param {string} signedChallenge - The signed challenge to verify.
   * @throws {Error} - Throws an error if the HTTP request fails.
   */
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

      window.location.reload();
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
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            id="username"
            type="text"
            placeholder="Enter your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <button onClick={handleGetSignedChallenge} disabled={loading}>
          {loading ? "Awaiting login" : "Login"}
        </button>
      </form>
      {message && <p className="message">{message}</p>}
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default LoginComponent;
