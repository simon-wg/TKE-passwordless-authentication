import React, { useState } from 'react';

const RegisterComponent = () => {
  const [username, setUsername] = useState('');
  const [pubkey, setPubkey] = useState('');

  const handleRegister = async () => {
    // Get the public key from the client server
    const pubKeyResponse = await fetch('http://localhost:8080/api/getTkeyPubKey');
    if (pubKeyResponse.ok) {
      const pubKeyData = await pubKeyResponse.json();
      setPubkey(pubKeyData.publicKey);
    } else {
      alert('Failed to get public key');
      return;
    }

    // Register the user with the public key
    const response = await fetch('http://localhost:8080/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, pubkey }),
    });

    if (response.ok) {
      alert('User registered successfully');
    } else {
      alert('Failed to register user');
    }
  };

  return (
    <div>
      <h2>Register</h2>
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <button onClick={handleRegister}>Register</button>
    </div>
  );
};

export default RegisterComponent;