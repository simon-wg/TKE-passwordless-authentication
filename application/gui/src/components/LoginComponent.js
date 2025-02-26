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
      body: JSON.stringify({ username }),
    });

    if (response.ok) {
      setSuccess('Successfully signed in')
    }

    else {
      setError('Failed to sign in user')
    }
  };

  const sign = async (challenge) => {
    // Sign the challenge using the client server
    const signResponse = await fetch(config.clientBaseUrl + '/api/sign', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ challenge }),
    });

    if (signResponse.ok) {
      const signData = await signResponse.json();
      console.log("Challenge signed");
      return signData.signature;
    } else {
      console.log("Failed to sign challenge");
      throw new Error('Failed to sign challenge');
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