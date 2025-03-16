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

  const handleRegister = async (event) => {
    event.preventDefault();
    setLoading(true);

    const result = await fetch("http://localhost:6060/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, label }),
    });

    if (result.ok) {
      setSuccess("Success!");
      setError("");
    } else {
      setSuccess("");
      setError("Error creating user");
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
