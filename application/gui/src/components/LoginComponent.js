import React, { useState } from 'react';
import './styles.css';
import config from '../config'

const LoginComponent = () => {
  const [username, setUsername] = useState('');
  const [message, setMessage] = useState('');
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');

  const handleLogin = async () => {
    // Get the challenge from the backend server
    setError('')
    const response = await fetch(config.serverBaseUrl + '/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username }),
    });

    if (response.ok) {
      const data = await response.json();
      console.log('Challenge received');

      // Sign the challenge using the client server
      setMessage('Touch TKey')
      let signedChallenge;
      try {
        signedChallenge = await sign(data.challenge);
        setMessage('')
      } catch (error) {
        setMessage('')
        setSuccess('');
        setError('Failed to sign challenge');
        return;
      }

      // Send signed challenge to application
      const submitResponse = await fetch(config.serverBaseUrl + '/api/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, signature: signedChallenge }),
      });

      if (submitResponse.ok) {
        setSuccess('Login successful');
        setError('');
      } else {
        setSuccess('');
        setError('Failed to submit signed challenge');
      }
    } else {
      setSuccess('');
      setError('Failed to login');
    }
  };

  const sign = async (challenge) => {
    // Sign the challenge using the client server
    const signResponse = await fetch(config.agentBaseUrl + '/api/sign', {
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
      {/* {challenge && (
        <div>
          <p>Challenge: {challenge}</p>
        </div>
      )} */}
      {/* {signature && <p>Signature: {signature}</p>} */}
      {message && <p className='message'>{message}</p>}
      {success && <p className="success">{success}</p>}
      {error && <p className="error">{error}</p>}
    </div>
  );
};

export default LoginComponent;