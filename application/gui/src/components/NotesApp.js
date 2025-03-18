import React, { useEffect, useState } from 'react';
import './styles.css';
import NoteCard from './NoteCard';
import './NotesApp.css';
import useFetchNotes from '../hooks/useFetchNotes';
import useAuthCheck from '../hooks/useAuthCheck';
import useDeleteNote from '../hooks/useDeleteNote';

const NotesApp = () => {
  const isAuthenticated = useAuthCheck();
  const fetchedNotes = useFetchNotes(isAuthenticated);
  const [notes, setNotes] = useState([]);
  const [selectedNote, setSelectedNotes] = useState(null);
  const [deleteResult, deleteNote] = useDeleteNote();

  useEffect(() => {
    setNotes(fetchedNotes);
  }, [fetchedNotes]);

  const handleNoteClick = (data) => {
    if (selectedNote && selectedNote.ID === data.ID) {
      setSelectedNotes(null);
    } else {
      setSelectedNotes(data);
    }
  };

  const handleAddNote = () => {
    const newNote = {
      ID: `temp-${Date.now()}`, // Temporary ID for unsaved notes
      Name: '',
      Note: '',
      isUnsaved: true,
    };
    console.log(newNote);
    console.log(notes);
    setNotes([...notes, newNote]);
    setSelectedNotes(newNote);
  };

  const handleUpdate = (updatedNote) => {
    console.log(notes);
    console.log(updatedNote);
    setNotes((prevNotes) =>
      prevNotes.map((note) =>
        note === selectedNote ? updatedNote : note
      )
    );
    setSelectedNotes(updatedNote);
  };

  const handleDelete = (id) => {
    // Check if the note is unsaved (temporary ID)
    if (id.startsWith('temp-')) {
      setNotes((prevNotes) =>
        prevNotes.filter((note) => note.ID !== id)
      );
      setSelectedNotes(null);
    } else {
      // Make HTTP request to delete saved note
      deleteNote(id).then(() => {
        setNotes((prevNotes) =>
          prevNotes.filter((note) => note.ID !== id)
        );
        setSelectedNotes(null);
      }).catch((error) => {
        console.error('Failed to delete note:', error);
      });
    }
  };

  return (
    <div className="note-manager">
      <div className="note-list">
        {notes.map((noteData) => (
          <div
            key={noteData.ID}
            className="note-list-item"
            onClick={() => handleNoteClick(noteData)}
          >
            {noteData.Name || 'New Note'}
          </div>
        ))}
        <button className="add-note-button" onClick={handleAddNote}>
          +
        </button>
      </div>
      <div className="note-details">
        {selectedNote && (
          <NoteCard
            key={selectedNote.ID}
            id={selectedNote.ID}
            name={selectedNote.Name}
            body={selectedNote.Password}
            isUnsaved={selectedNote.isUnsaved || false}
            onUpdate={handleUpdate}
            onDelete={handleDelete}
          />
        )}
      </div>
    </div>
  );
};

export default NotesApp;