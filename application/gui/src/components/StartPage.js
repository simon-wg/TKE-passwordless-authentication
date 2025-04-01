import React from 'react';
import './StartPage.css';

const StartPage = ({ setPage }) => {
    const handleRegisterClick = () => {
        setPage("register");
    };

    return (
        <div style={{ textAlign: 'center', padding: '20px' }}>
            <h1>Welcome to the TKey Login Demo</h1>
            <p>
                This website demonstrates a passwordless authentication system using TKey. 
            </p>
            <p>
                Make sure the client service is running on your device before continuing.
            </p>
            <p>
                When the client is running, plug in your TKey, head to the register or login page and fill out the form as usual.
            </p>
            <p>
                The first time registering or login in after plugging in you will be prompted by an optional user supplied secret (password)."
            </p>
            <p>
                Then when the TKey starts to blink green, give it a touch!
            </p>
            <button className='to-register-btn' onClick={handleRegisterClick}>
                Register with TKey
            </button>
        </div>
    );
};

export default StartPage;