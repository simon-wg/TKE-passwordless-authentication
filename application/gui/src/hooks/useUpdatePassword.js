import { useState } from 'react';

const useUpdatePassword = () => {
  const [result, setResult] = useState(null);

  const updatePassword = async (id, name, password, isAuthenticated) => {
    if (!isAuthenticated) {
      setResult(false);
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/api/update-password', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id, name, password }),
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

  return [result, updatePassword];
};

export default useUpdatePassword;