import React, { useState } from 'react';
import './styles.css';
import PasswordCard from './PasswordCard';
import './PasswordApp.css';
import useFetchPasswords from '../hooks/useFetchPasswords'

const PasswordApp = () => {
    const passwords = [
        { name: 'First Note', password: 'This is the body of the first note.' },
        { name: 'Second Note', password: 'This is the body of the second note.' }
    ].map((note, index) => ({ ...note, id: index + 1 }));
    const [selectedNote, setSelectedPassword] = useState(null);
    const userPasswords = useFetchPasswords(true);


    const handleNoteClick = (data) => {
        if (selectedNote && selectedNote.id === data.id) {
            setSelectedPassword(null);
        } else {
            setSelectedPassword(data);
        }
    };

  return (
    <div className="password-manager">
        <div className="password-list">
            {passwords.map((passwordData) => (
            <div
                key={passwordData.id}
                className="password-list-item"
                onClick={() => handleNoteClick(passwordData)}
            >
                {passwordData.name}
            </div>
            ))}
        </div>
        <div className="password-details">
            {selectedNote && (
            <PasswordCard
                key={selectedNote.id}
                name={selectedNote.name}
                body={selectedNote.password}
            />
            )}
        </div>
    </div>
    );
};

export default PasswordApp;