import React, { useState } from 'react';
import './styles.css';
import config from '../config'

const RegisterComponent = () => {
  const [username, setUsername] = useState('');
  const [pubkey, setPubkey] = useState('');
  const [success, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleRegister = async (event) => {
    event.preventDefault(); // Prevent the default form submission behavior

    // Get the public key from the client server
    const pubKeyResponse = await fetch(config.clientBaseUrl + '/api/getTkeyPubKey');
    if (pubKeyResponse.ok) {
      const pubKeyData = await pubKeyResponse.json();
      console.log("Pubkey: " + pubKeyData.publicKey);
      setPubkey(pubKeyData.publicKey);

      // Register the user with the public key
      const response = await fetch(config.serverBaseUrl + '/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, pubkey: pubKeyData.publicKey }),
      });

      if (response.ok) {
        setMessage('User registered successfully');
        setError('');
      } else {
        setMessage('');
        setError('Failed to register user');
      }
    } else {
      setMessage('');
      setError('Failed to get public key');
      return;
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