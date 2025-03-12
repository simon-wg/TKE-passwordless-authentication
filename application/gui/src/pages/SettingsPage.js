import React, { useState, useEffect } from "react";
import useAuthCheck from "../hooks/useAuthCheck";
import useFetchUser from "../hooks/useFetchUser";
import "../components/styles.css";

const SettingsPage = () => {
  const isAuthenticated = useAuthCheck();
  const user = useFetchUser(isAuthenticated);
  const [addKeyLabel, setAddKeyLabel] = useState("");
  const [removeKeyLabel, setRemoveKeyLabel] = useState("");
  const [message, setMessage] = useState("");
  const [keyLabels, setKeyLabels] = useState([]);

  useEffect(() => {
    const fetchKeyLabels = async () => {
      const response = await fetch(
        "http://localhost:8080/api/get-public-key-labels",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ username: user }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        setKeyLabels(data.labels);
      } else {
        setMessage("Error fetching public key labels");
      }
    };

    if (user) {
      fetchKeyLabels();
    }
  }, [user]);

  const handleAddKey = async () => {
    const response = await fetch("http://localhost:6060/api/add-public-key", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username: user, label: addKeyLabel }),
    });

    if (response.ok) {
      setMessage("Public key added successfully");
      setAddKeyLabel("");
      // Refresh the key labels
      const data = await response.json();
      setKeyLabels(data.labels);
    } else {
      setMessage("Error adding public key");
    }
  };

  const handleRemoveKey = async () => {
    const response = await fetch(
      "http://localhost:6060/api/remove-public-key",
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
      setRemoveKeyLabel("");
      // Refresh the key labels
      const data = await response.json();
      setKeyLabels(data.labels);
    } else {
      setMessage("Error removing public key");
    }
  };

  if (!isAuthenticated) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container">
      <h1>Settings</h1>
      <div>
        <h2>Public Key Labels</h2>
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
      {message && <p>{message}</p>}
    </div>
  );
};

export default SettingsPage;
