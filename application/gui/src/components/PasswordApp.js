import React, { useEffect, useState } from 'react';
import './styles.css';
import PasswordCard from './PasswordCard';
import './PasswordApp.css';
import useFetchPasswords from '../hooks/useFetchPasswords';
import useAuthCheck from '../hooks/useAuthCheck';
import useDeletePassword from '../hooks/useDeletePassword';

const PasswordApp = () => {
  const isAuthenticated = useAuthCheck();
  const fetchedPasswords = useFetchPasswords(isAuthenticated);
  const [passwords, setPasswords] = useState([]);
  const [selectedNote, setSelectedPassword] = useState(null);
  const [deleteResult, deletePassword] = useDeletePassword();

  useEffect(() => {
    setPasswords(fetchedPasswords);
  }, [fetchedPasswords]);

  const handleNoteClick = (data) => {
    if (selectedNote && selectedNote.ID === data.ID) {
      setSelectedPassword(null);
    } else {
      setSelectedPassword(data);
    }
  };

  const handleAddPassword = () => {
    const newPassword = {
      ID: `temp-${Date.now()}`, // Temporary ID for unsaved passwords
      Name: '',
      Password: '',
      isUnsaved: true,
    };
    console.log(newPassword);
    setPasswords([...passwords, newPassword]);
    setSelectedPassword(newPassword);
  };

  const handleUpdate = (updatedPassword) => {
    console.log(passwords);
    console.log(updatedPassword);
    setPasswords((prevPasswords) =>
      prevPasswords.map((password) =>
        password === selectedNote ? updatedPassword : password
      )
    );
    setSelectedPassword(updatedPassword);
  };

  const handleDelete = (id) => {
    // Check if the password is unsaved (temporary ID)
    if (id.startsWith('temp-')) {
      setPasswords((prevPasswords) =>
        prevPasswords.filter((password) => password.ID !== id)
      );
      setSelectedPassword(null);
    } else {
      // Make HTTP request to delete saved password
      deletePassword(id).then(() => {
        setPasswords((prevPasswords) =>
          prevPasswords.filter((password) => password.ID !== id)
        );
        setSelectedPassword(null);
      }).catch((error) => {
        console.error('Failed to delete password:', error);
      });
    }
  };
  

  return (
    <div className="password-manager">
      <div className="password-list">
        {passwords.map((passwordData) => (
          <div
            key={passwordData.ID}
            className="password-list-item"
            onClick={() => handleNoteClick(passwordData)}
          >
            {passwordData.Name || 'New Password'}
          </div>
        ))}
        <button className="add-password-button" onClick={handleAddPassword}>
          +
        </button>
      </div>
      <div className="password-details">
        {selectedNote && (
          <PasswordCard
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

export default PasswordApp;