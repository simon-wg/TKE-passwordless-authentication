import React, { useState } from 'react';
import './styles.css';
import PasswordCard from './PasswordCard';
import './PasswordApp.css';
import useFetchPasswords from '../hooks/useFetchPasswords';
import useAuthCheck from '../hooks/useAuthCheck';

const PasswordApp = () => {
  const isAuthenticated = useAuthCheck();
  const userPasswords = useFetchPasswords(isAuthenticated);
  const [selectedNote, setSelectedPassword] = useState(null);

  const handleNoteClick = (data) => {
    if (selectedNote && selectedNote.ID === data.id) {
      setSelectedPassword(null);
    } else {
      setSelectedPassword(data);
    }
  };

  return (
    <div className="password-manager">
      <div className="password-list">
        {userPasswords.map((passwordData) => (
          <div
            key={passwordData.ID}
            className="password-list-item"
            onClick={() => handleNoteClick(passwordData)}
          >
            {passwordData.Name}
          </div>
        ))}
      </div>
      <div className="password-details">
        {selectedNote && (
          <PasswordCard
            key={selectedNote.ID}
            name={selectedNote.Name}
            body={selectedNote.Password}
          />
        )}
      </div>
    </div>
  );
};

export default PasswordApp;