import React, { useState, useEffect } from 'react';
import useAuthCheck from "../hooks/useAuthCheck";
import useSavePassword from '../hooks/useSavePassword';
import './PasswordCard.css';

const PasswordCard = ({ name: initialName, body: initialBody, isNew }) => {
  const [name, setName] = useState(initialName);
  const [body, setBody] = useState(initialBody);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState(''); 
  const [saveClicked, setSaveClicked] = useState(false);

  const isAuthenticated = useAuthCheck();
  const endpoint = 'save-password';
  const [result, savePassword] = useSavePassword();

  useEffect(() => {
    if (saveClicked) {
      savePassword(name, body, isAuthenticated, endpoint);
      setSaveClicked(false);
    }
  }, [saveClicked, name, body, isAuthenticated, endpoint, savePassword]);

  useEffect(() => {
    if (result !== null) {
      if (result === false) {
        setMessage('Failed to save password');
        setMessageType('error');
      } else {
        setMessage('Password saved successfully');
        setMessageType('success');
      }
    }
  }, [result]);

  const handleNameChange = (event) => {
    setName(event.target.value);
  };

  const handleBodyChange = (event) => {
    setBody(event.target.value);
  };

  const handleSaveClick = (event) => {
    event.preventDefault();
    setSaveClicked(true);
  };

  return (
    <div className="password-card">
      <input
        type="text"
        value={name}
        onChange={handleNameChange}
        placeholder="Name"
      />
      <textarea
        value={body}
        onChange={handleBodyChange}
        placeholder="Body"
      />
      <button onClick={handleSaveClick}>Save</button>
      {message && <p className={messageType}>{message}</p>}
    </div>
  );
};

export default PasswordCard;