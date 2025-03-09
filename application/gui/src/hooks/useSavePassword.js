import { useState } from 'react';

const useSavePassword = () => {
  const [result, setResult] = useState(null);

  const savePassword = async (name, password, isAuthenticated, endpoint) => {
    if (!isAuthenticated) {
      setResult(false);
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/api/' + endpoint, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, password }),
      });

      if (response.ok) {
        setResult(true);
      } else {
        setResult(false);
      }
    } catch (error) {
      console.log('Error saving password', error);
      setResult(false);
    }
  };

  return [result, savePassword];
};

export default useSavePassword;