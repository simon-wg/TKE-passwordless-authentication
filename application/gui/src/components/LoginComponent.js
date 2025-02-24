// import React, { useState } from 'react';

// const LoginComponent = () => {
//   const [username, setUsername] = useState('');
//   const [challenge, setChallenge] = useState('');
//   const [signature, setSignature] = useState('');

//   const handleLogin = async () => {
//     // Get the challenge from the backend server
//     const response = await fetch('http://192.168.50.106:8080/api/login', {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify({ username }),
//     });

//     if (response.ok) {
//       const data = await response.json();
//       setChallenge(data.challenge);
//       alert('Challenge received');
//     } else {
//       alert('Failed to login');
//     }
//   };

//   const handleSign = async () => {
//     // Sign the challenge using the client server
//     const signResponse = await fetch('http://localhost:8080/api/sign', {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify({ challenge }),
//     });

//     if (signResponse.ok) {
//         const signData = await signResponse.json();
//         setSignature(signData.signature);
//         console.log("Challenge signed")

//         // Send signed challenge to application
//         const submitResponse = await fetch('http://192.168.50.106:8080/api/verify', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({ username, signature }),
//         });

//         if (submitResponse.ok) {
//             alert('Login successful');
//         } else {
//             alert('Failed to submit signed challenge');
//         }
//     }
//     else {
//         console.log("Failed to sign challenge")
//     }
//   };

//   return (
//     <div>
//       <h2>Login</h2>
//       <input
//         type="text"
//         placeholder="Username"
//         value={username}
//         onChange={(e) => setUsername(e.target.value)}
//       />
//       <button onClick={handleLogin}>Login</button>
//       {challenge && (
//         <div>
//           <p>Challenge: {challenge}</p>
//           <button onClick={handleSign}>Sign Challenge</button>
//         </div>
//       )}
//       {signature && <p>Signature: {signature}</p>}
//     </div>
//   );
// };

// export default LoginComponent;

import React, { useState } from 'react';

const LoginComponent = () => {
  const [username, setUsername] = useState('');
  const [challenge, setChallenge] = useState('');
  const [signature, setSignature] = useState('');

  const handleLogin = async () => {
    // Get the challenge from the backend server
    const response = await fetch('http://192.168.50.106:8080/api/login', {
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

      // Sign the challenge using the client server
      const signedChallenge = await sign(data.challenge);

      // Send signed challenge to application
      const submitResponse = await fetch('http://192.168.50.106:8080/api/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, signature: signedChallenge }),
      });

      if (submitResponse.ok) {
        alert('Login successful');
      } else {
        alert('Failed to submit signed challenge');
      }
    } else {
      alert('Failed to login');
    }
  };

  const sign = async (challenge) => {
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
      console.log("Challenge signed");
      return signData.signature;
    } else {
      console.log("Failed to sign challenge");
      throw new Error('Failed to sign challenge');
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
        </div>
      )}
      {signature && <p>Signature: {signature}</p>}
    </div>
  );
};

export default LoginComponent;