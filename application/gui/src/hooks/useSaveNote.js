import { useState } from 'react';

const useSaveNote = () => {
  const [result, setResult] = useState(null);

  const saveNote = async (name, note, isAuthenticated, endpoint) => {
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
        body: JSON.stringify({ name, password: note }),
      });

      if (response.ok) {
        const data = await response.json();
        setResult(data);
      } else {
        setResult(false);
      }
    } catch (error) {
      console.log('Error saving note', error);
      setResult(false);
    }
  };

  return [result, saveNote];
};

export default useSaveNote;