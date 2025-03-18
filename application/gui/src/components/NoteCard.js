import React, { useState, useEffect, useRef } from 'react';
import useCreateNote from '../hooks/useSaveNote';
import useUpdateNote from '../hooks/useUpdateNote';
import useDeleteNote from '../hooks/useDeleteNote';
import './NoteCard.css';

const NoteCard = ({ id: initialId, name: initialName, body: initialBody, isUnsaved: unsavedInitial = false, onUpdate, onDelete }) => {
  const [id, setId] = useState(initialId);
  const [name, setName] = useState(initialName);
  const [body, setBody] = useState(initialBody);
  const [isUnsaved, setIsNew] = useState(unsavedInitial);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('');
  const [saveClicked, setSaveClicked] = useState(false);

  const [saveResult, saveNote] = useCreateNote();
  const [updateResult, updateNote] = useUpdateNote();
  const [deleteResult, deleteNote] = useDeleteNote();

  const prevSaveResult = useRef(null);
  const prevUpdateResult = useRef(null);
  const prevDeleteResult = useRef(null);

  useEffect(() => {
    if (!saveClicked) return;
    if (isUnsaved) {
      saveNote(name, body);
      setIsNew(false);
    } else {
      updateNote(id, name, body);
    }
    setSaveClicked(false);
  }, [saveClicked, saveNote, updateNote, isUnsaved, name, body, id]);

  useEffect(() => {
    if (saveResult !== null && saveResult !== prevSaveResult.current) {
      if (saveResult === false) {
        setMessage('Failed to create note');
        setMessageType('error');
      } else {
        setMessage('Note created successfully');
        setMessageType('success');
        setId(saveResult.id);
        if (onUpdate) onUpdate({ ID: saveResult.id, Name: name, Note: body });
      }
      prevSaveResult.current = saveResult;
    }
  }, [saveResult, onUpdate]);

  useEffect(() => {
    console.log(message);
  }, [message])

  useEffect(() => {
    if (updateResult !== null && updateResult !== prevUpdateResult.current) {
      if (updateResult === false) {
        setMessage('Failed to update note');
        setMessageType('error');
      } else {
        setMessage('Note updated successfully');
        setMessageType('success');
        if (onUpdate) onUpdate({ ID: id, Name: name, Note: body });
      }
      prevUpdateResult.current = updateResult;
    }
  }, [updateResult, onUpdate, id, name, body]);

  useEffect(() => {
    if (deleteResult !== null && deleteResult !== prevDeleteResult.current) {
      if (deleteResult === false) {
        setMessage('Failed to delete note');
        setMessageType('error');
      } else {
        setMessage('Note deleted successfully');
        setMessageType('success');
        if (onDelete) onDelete(id);
      }
      prevDeleteResult.current = deleteResult;
    }
  }, [deleteResult, onDelete, id]);

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
    event.preventDefault();
    onDelete(id);
  };

  return (
    <div className="note-card">
      <input
        type="text"
        value={name}
        onChange={handleNameChange}
        placeholder="Name"
      />
      <textarea
        value={body}
        onChange={handleBodyChange}
        placeholder="Note"
      />
      <div className="button-group">
        <button onClick={handleSaveClick}>Save</button>
        <button className="delete" onClick={handleDeleteClick}>Delete</button>
      </div>
      {message && <p className={messageType}>{message}</p>}
    </div>
  );
};

export default NoteCard;