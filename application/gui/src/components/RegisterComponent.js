import React, { useState } from "react";
import "./styles.css";
import config from "../config";
import LoadingCircle from "./LoadingCircle";
const RegisterComponent = () => {
  const [username, setUsername] = useState("");
  const [label, setLabel] = useState("");
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  /**
   * handleRegister sends a registration request to the client base URL with the provided username and label.
   * It sets the loading state while the request is being processed and updates the success or error state based on the response.
   *
   * Parameters:
   * - event: The event object associated with the form submission
   *
   * Returns:
   * - None
   */
  const handleRegister = async (event) => {
    event.preventDefault();
    setLoading(true);

    // Sends POST register request to the client.
    const result = await fetch(config.clientBaseUrl + "/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, label }),
    });

    // Checks whether fetch response was successful or not. And responds accordingly.
    if (result.ok) {
      setSuccess("Success!");
      setError("");
    } else {

      // Retrieves potential error message retrieved from the http error response and displays it to the user.
      var errorMessage = await result.text();
      setSuccess("");
      setError(errorMessage);
    }
    setLoading(false);
  };

  return (
    <div className="container">
      <h2>Register</h2>
      <LoadingCircle loading={loading} />

      <form onSubmit={handleRegister}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="text"
          placeholder="Key Label"
          value={label}
          onChange={(e) => setLabel(e.target.value)}
        />
        <button onClick={handleRegister} disabled={loading}>
          {loading ? "Loading..." : "Register"}
        </button>
      </form>
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default RegisterComponent;
