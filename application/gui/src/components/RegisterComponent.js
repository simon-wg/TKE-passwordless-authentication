import React, { useState } from 'react';
import './styles.css';
import config from '../config'

const RegisterComponent = () => {
  const [username, setUsername] = useState('');
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');

  const handleRegister = async (event) => {
    event.preventDefault();

    const result = await fetch(config.clientBaseUrl + '/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username }),
    })

    if (result.ok) {
      setSuccess('Success!')
      setError('')
    }

    else {
      setSuccess('')
      setError('Error creating user')
    }
  };

  return (
    <div className="container">
      <h2>Register</h2>
      <form onSubmit={handleRegister}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button type="submit">Register</button>
      </form>
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default RegisterComponent;