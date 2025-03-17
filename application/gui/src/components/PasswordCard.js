import React, { useState, useEffect, useRef } from 'react';
import useAuthCheck from "../hooks/useAuthCheck";
import useSavePassword from '../hooks/useSavePassword';
import useUpdatePassword from '../hooks/useUpdatePassword';
import useDeletePassword from '../hooks/useDeletePassword';
import './PasswordCard.css';

const PasswordCard = ({ id: initialId, name: initialName, body: initialBody, isUnsaved : unsavedInitial = false, onUpdate, onDelete }) => {
  const [id, setId] = useState(initialId);
  const [name, setName] = useState(initialName);
  const [body, setBody] = useState(initialBody);
  const [isUnsaved, setIsNew] = useState(unsavedInitial);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState(''); 
  const [saveClicked, setSaveClicked] = useState(false);

  const isAuthenticated = useAuthCheck();
  const [saveResult, savePassword] = useSavePassword();
  const [updateResult, updatePassword] = useUpdatePassword();
  const [deleteResult, deletePassword] = useDeletePassword();

  const prevSaveResult = useRef(null);
  const prevUpdateResult = useRef(null);
  const prevDeleteResult = useRef(null);

  useEffect(() => {
    if (!saveClicked) return;
    if (isUnsaved) {
      savePassword(name, body, isAuthenticated, 'save-password');
      setIsNew(false);
    } else {
      updatePassword(id, name, body, isAuthenticated, 'update-password');
    }
    setSaveClicked(false);
  }, [saveClicked, savePassword, updatePassword, isUnsaved, name, body, isAuthenticated, id]);

  useEffect(() => {
    if (saveResult !== null && saveResult !== prevSaveResult.current) {
      if (saveResult === false) {
        setMessage('Failed to create password');
        setMessageType('error');
      } else {
        setMessage('Password created successfully');
        setMessageType('success');
        setId(saveResult.id);
        if (onUpdate) onUpdate({ ID: saveResult.id, Name: name, Password: body });
      }
      prevSaveResult.current = saveResult;
    }
  }, [saveResult, onUpdate, name, body]);

  useEffect(() => {
    if (updateResult !== null && updateResult !== prevUpdateResult.current) {
      if (updateResult === false) {
        setMessage('Failed to update password');
        setMessageType('error');
      } else {
        setMessage('Password updated successfully');
        setMessageType('success');
        if (onUpdate) onUpdate({ ID: id, Name: name, Password: body });
      }
      prevUpdateResult.current = updateResult;
    }
  }, [updateResult, onUpdate, id, name, body]);

  useEffect(() => {
    if (deleteResult !== null && deleteResult !== prevDeleteResult.current) {
      if (deleteResult === false) {
        setMessage('Failed to delete password');
        setMessageType('error');
      } else {
        setMessage('Password deleted successfully');
        setMessageType('success');
        if (onDelete) onDelete(id);
      }
      prevDeleteResult.current = deleteResult;
    }
  }, [deleteResult, onDelete]);

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

  const handleDeleteClick = (event) => {
    if (isUnsaved) {
      onDelete(id);
    }
    else {
      event.preventDefault();
      deletePassword(id);
    }
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
      <div className="button-group">
        <button onClick={handleSaveClick}>Save</button>
        <button className="delete" onClick={handleDeleteClick}>Delete</button>
      </div>
      {message && <p className={messageType}>{message}</p>}
    </div>
  );
};

export default PasswordCard;