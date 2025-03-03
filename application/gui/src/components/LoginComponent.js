import React, { useState } from 'react';
import './styles.css';
import config from '../config'

const LoginComponent = () => {
  const [username, setUsername] = useState('');
  const [message, setMessage] = useState('');
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');

  const handleLogin = async () => {
    clearMessages()

    const response = await fetch(config.clientBaseUrl + '/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ username }),
    });

    if (response.ok) {
      setSuccess('Successfully signed in')
    }

    else {
      setError('Failed to sign in user')
    }
  };

  const clearMessages = async () => {
    setMessage('');
    setSuccess('');
    setError('');
  }

  return (
    <div className="container">
      <h2>Login</h2>
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <button onClick={handleLogin}>Login</button>
      {message && <p className='message'>{message}</p>}
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default LoginComponent;