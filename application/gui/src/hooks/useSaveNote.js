import { useState } from 'react';

const useCreateNote = () => {
  const [saveResult, setSagveResult] = useState(null);

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
        setSagveResult(data);
      } else {
        setSagveResult(false);
      }
    } catch (error) {
      console.log('Error creating note', error);
      setSagveResult(false);
    }
  };

  return [saveResult, createNote];
};

export default useCreateNote;