import React, { useState, useEffect } from "react";
import useFetchUser from "../hooks/useFetchUser";
import config from "../config";
import "../components/styles.css";
import { useNavigate } from "react-router-dom";




const SettingsPage = () => {
  const navigate = useNavigate();
  const user = useFetchUser();
  const [addKeyLabel, setAddKeyLabel] = useState("");
  const [removeKeyLabel, setRemoveKeyLabel] = useState("");
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");
  const [keyLabels, setKeyLabels] = useState([]);
  const [showDeletePopup, setShowDeletePopup] = useState(false);
  const [deleteConfirmation, setDeleteConfirmation] = useState("");
  const [popupMessage, setPopupMessage] = useState("");

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
      body: JSON.stringify({ label: addKeyLabel }),
    });

    if (response.ok) {
      setMessage("Public key added successfully");
      setMessageType("success");
      setAddKeyLabel("");
      fetchKeyLabels();
    } else {
      setMessage("Error adding public key");
      setMessageType("error");
    }
  };

  const handleRemoveKey = async () => {
    const response = await fetch(config.clientBaseUrl + "/api/remove-public-key", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ label: removeKeyLabel }),
    });

    if (response.ok) {
      setMessage("Public key removed successfully");
      setMessageType("success");
      setRemoveKeyLabel("");
      fetchKeyLabels();
    } else {
      setMessage("Error removing public key");
      setMessageType("error");
    }
  };

  const handleAccountDeletion = async () => {
    if (deleteConfirmation === "REMOVEMYACCOUNT") {
      const response = await fetch("/api/unregister", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });

      if (response.ok) {
        setMessage("Account deleted successfully");
        setMessageType("success");
        setShowDeletePopup(false);

        // !! TEMPORARY SOLUTION. api/logout should be called here when implemented. !!
        navigate('/');
      } else {
        setPopupMessage("Error deleting account");
      }
    } else {
      setPopupMessage("Incorrect confirmation input");
    }
  };

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
      <button className="delete-button" onClick={() => setShowDeletePopup(true)}>
        Delete Account
      </button>

      {showDeletePopup && (
        <div className="lightbox">
          <div className="popup">
            <h2>Confirm Account Deletion</h2>
            <p>Type "REMOVEMYACCOUNT" to confirm:</p>
            <input
              type="text"
              value={deleteConfirmation}
              onChange={(e) => setDeleteConfirmation(e.target.value)}
            />
            <button style={{ marginBottom: "10px" }} onClick={handleAccountDeletion}>Confirm</button>
            <button onClick={() => setShowDeletePopup(false)}>Cancel</button>
            {popupMessage && <p className="error-message" style={{ color: "red" }}>{popupMessage}</p>}
          </div>
        </div>
      )}
    </div>
  );
};

export default SettingsPage;
