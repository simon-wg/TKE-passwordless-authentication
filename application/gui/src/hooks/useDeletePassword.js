import { useState } from 'react';
import useAuthCheck from './useAuthCheck';

const useDeletePassword = () => {
  const [result, setResult] = useState(null);
  const isAuthenticated = useAuthCheck();

  const deletePassword = async (id) => {
    if (!isAuthenticated) return;

    try {
      const response = await fetch('http://localhost:8080/api/delete-password', {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id }),
      });

      if (response.ok) {
        const data = await response.json();
        setResult(data);
      } else {
        setResult(false);
      }
    } catch (error) {
      console.log('Error deleting password', error);
      setResult(false);
    }
  };

  return [result, deletePassword];
};

export default useDeletePassword;