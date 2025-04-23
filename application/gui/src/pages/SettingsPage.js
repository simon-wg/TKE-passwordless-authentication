import React, { useState, useEffect } from "react";
import useFetchUser from "../hooks/useFetchUser";
import config from "../config";
import "../components/styles.css";
import { secureFetch } from "../util/secureFetch";
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
    const response = await secureFetch("/api/get-public-key-labels", {
      method: "POST",
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
    try {
      const clientResponse = await fetch(
        config.clientBaseUrl + "/api/add-public-key",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (!clientResponse.ok) {
        const errorText = await clientResponse.text();
        setMessage(errorText);
        setMessageType("error");
        return;
      }

      const { pubkey } = await clientResponse.json();

      const backendResponse = await secureFetch("/api/add-public-key", {
        method: "POST",
        body: JSON.stringify({ label: addKeyLabel, pubkey }),
      });

      if (!backendResponse.ok) {
        const errorText = await backendResponse.text();
        setMessage(errorText);
        setMessageType("error");
        return;
      }

      setMessage("Public key added successfully");
      setMessageType("success");
      setAddKeyLabel("");
      fetchKeyLabels();
    } catch (error) {
      console.error("Error in add key flow:", error);
      setMessage(`Unexpected Error: ${error.message}`);
      setMessageType("error");
    }
  };

  const handleRemoveKey = async () => {
    const response = await secureFetch("/api/remove-public-key", {
      method: "POST",
      body: JSON.stringify({ label: removeKeyLabel }),
    });

    if (response.ok) {
      setMessage("Public key removed successfully");
      setMessageType("success");
      setRemoveKeyLabel("");
      fetchKeyLabels();
    } else {
      const errorText = await response.text();
      setMessage(errorText);
      setMessageType("error");
      return;
    }
  };

  const handleAccountDeletion = async () => {
    if (deleteConfirmation === "REMOVEMYACCOUNT") {
      const response = await secureFetch("/api/unregister", {
        method: "POST",
      });

      if (response.ok) {
        setMessage("Account deleted successfully");
        setMessageType("success");
        setShowDeletePopup(false);
        navigate("/");
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

      {/* Message Display */}
      {message && <p className={`message ${messageType}`}>{message}</p>}

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
        <div className="form-group">
          <label htmlFor="addKeyLabel">Key Label</label>
          <input
            id="addKeyLabel"
            type="text"
            placeholder="Enter a label for the new key"
            value={addKeyLabel}
            onChange={(e) => setAddKeyLabel(e.target.value)}
          />
        </div>
        <button onClick={handleAddKey}>Add Public Key</button>
      </div>
      <div>
        <h2>Remove Public Key</h2>
        <div className="form-group">
          <label htmlFor="removeKeyLabel">Key Label</label>
          <input
            id="removeKeyLabel"
            type="text"
            placeholder="Enter the label of the key to remove"
            value={removeKeyLabel}
            onChange={(e) => setRemoveKeyLabel(e.target.value)}
          />
        </div>
        <button onClick={handleRemoveKey}>Remove Public Key</button>
      </div>

      <div>
        <h2>Account Deletion</h2>
        <button
          className="delete-button"
          onClick={() => setShowDeletePopup(true)}
        >
          Delete Account
        </button>
      </div>

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
            <button
              className="delete-button"
              style={{ marginBottom: "10px" }}
              onClick={handleAccountDeletion}
            >
              Confirm
            </button>
            <button onClick={() => setShowDeletePopup(false)}>Cancel</button>
            {popupMessage && (
              <p className="error-message" style={{ color: "red" }}>
                {popupMessage}
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default SettingsPage;
