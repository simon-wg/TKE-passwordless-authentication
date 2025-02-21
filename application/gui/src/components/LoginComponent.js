import React, { useState } from 'react';

const LoginComponent = () => {
  const [username, setUsername] = useState('');
  const [challenge, setChallenge] = useState('');
  const [signature, setSignature] = useState('');

  const handleLogin = async () => {
    // Get the challenge from the backend server
    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username }),
    });

    if (response.ok) {
      const data = await response.json();
      setChallenge(data.challenge);
      alert('Challenge received');
    } else {
      alert('Failed to login');
    }
  };

  const handleSign = async () => {
    // Sign the challenge using the client server
    const signResponse = await fetch('http://localhost:8080/api/sign', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ challenge }),
    });

    if (signResponse.ok) {
      const signData = await signResponse.json();
      setSignature(signData.signature);
      alert('Challenge signed');
    } else {
      alert('Failed to sign challenge');
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <button onClick={handleLogin}>Login</button>
      {challenge && (
        <div>
          <p>Challenge: {challenge}</p>
          <button onClick={handleSign}>Sign Challenge</button>
        </div>
      )}
      {signature && <p>Signature: {signature}</p>}
    </div>
  );
};

export default LoginComponent;