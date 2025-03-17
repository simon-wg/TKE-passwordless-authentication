import React, { useState, useEffect, useNavigate } from "react";
import useFetchUser from "../hooks/useFetchUser";
import config from "../config";
import "../components/styles.css";

const SettingsPage = () => {
  const user = useFetchUser();
  const [addKeyLabel, setAddKeyLabel] = useState("");
  const [removeKeyLabel, setRemoveKeyLabel] = useState("");
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");
  const [keyLabels, setKeyLabels] = useState([]);

  const navigate = useNavigate();

  const fetchKeyLabels = async () => {
    const response = await fetch("/api/get-public-key-labels", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username: user }),
    });

    if (response.ok) {
      const data = await response.json();
      setKeyLabels(data.labels);
    } else {
      setMessage("Error fetching public key labels");
      setMessageType("error");
    }
  };

  useEffect(() => {
    if (user) {
      fetchKeyLabels();
    }
  }, [user]);

  const handleAddKey = async () => {
    const response = await fetch(config.clientBaseUrl + "/api/add-public-key", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username: user, label: addKeyLabel }),
    });

    if (response.ok) {
      setMessage("Public key added successfully");
      setMessageType("success");
      setAddKeyLabel("");
      // Refresh the key labels
      fetchKeyLabels();
    } else {
      setMessage("Error adding public key");
      setMessageType("error");
    }
  };

  const handleRemoveKey = async () => {
    const response = await fetch(
      config.clientBaseUrl + "/api/remove-public-key",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ username: user, label: removeKeyLabel }),
      }
    );

    if (response.ok) {
      setMessage("Public key removed successfully");
      setMessageType("success");
      setRemoveKeyLabel("");
      // Refresh the key labels
      fetchKeyLabels();
    } else {
      setMessage("Error removing public key");
      setMessageType("error");
    }
  };

  if (!user) {
    navigate("/");
  }

  return (
    <div className="container">
      <h1>Settings</h1>
      <div>
        <h2>Your Public Keys</h2>
        <ul>
          {keyLabels.map((label, index) => (
            <li key={index}>{label}</li>
          ))}
        </ul>
      </div>
      <div>
        <h2>Add Public Key</h2>
        <input
          type="text"
          placeholder="Key Label"
          value={addKeyLabel}
          onChange={(e) => setAddKeyLabel(e.target.value)}
        />
        <button onClick={handleAddKey}>Add Public Key</button>
      </div>
      <div>
        <h2>Remove Public Key</h2>
        <input
          type="text"
          placeholder="Key Label"
          value={removeKeyLabel}
          onChange={(e) => setRemoveKeyLabel(e.target.value)}
        />
        <button onClick={handleRemoveKey}>Remove Public Key</button>
      </div>
      {message && <p className={messageType}>{message}</p>}
    </div>
  );
};

export default SettingsPage;
