import React, { useEffect, useState } from 'react';
import './styles.css';
import PasswordCard from './PasswordCard';
import './PasswordApp.css';
import useFetchPasswords from '../hooks/useFetchPasswords';
import useAuthCheck from '../hooks/useAuthCheck';

const PasswordApp = () => {
  const isAuthenticated = useAuthCheck();
  const fetchedPasswords = useFetchPasswords(isAuthenticated);
  const [passwords, setPasswords] = useState([]);
  const [selectedNote, setSelectedPassword] = useState(null);

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
      ID: passwords.length + 1,
      Name: '',
      Password: '',
      isNewPassword: true,
    };
    setPasswords([...passwords, newPassword]);
    setSelectedPassword(newPassword);
  };

  const handleUpdate = (updatedPassword) => {
    setPasswords((prevPasswords) =>
      prevPasswords.map((password) =>
        password.ID === updatedPassword.ID ? updatedPassword : password
      )
    );
    setSelectedPassword(updatedPassword);
  };

  const handleDelete = (id) => {
    setPasswords((prevPasswords) =>
      prevPasswords.filter((password) => password.ID !== id)
    );
    setSelectedPassword(null);
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
            isNewPassword={selectedNote.isNewPassword || false}
            onUpdate={handleUpdate}
            onDelete={handleDelete}
          />
        )}
      </div>
    </div>
  );
};

export default PasswordApp;