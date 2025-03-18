import { useState } from 'react';

const useCreateNote = () => {
  const [result, setResult] = useState(null);

  const createNote = async (name, note) => {
    try {
      const response = await fetch('/api/create-note', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, note }),
      });

      if (response.ok) {
        const data = await response.json();
        setResult(data);
      } else {
        setResult(false);
      }
    } catch (error) {
      console.log('Error creating note', error);
      setResult(false);
    }
  };

  return [result, createNote];
};

export default useCreateNote;