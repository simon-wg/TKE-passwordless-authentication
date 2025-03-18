import { useState } from 'react';

const useUpdateNote = () => {
  const [result, setResult] = useState(null);

  const updateNote = async (id, name, note) => {
    try {
      const response = await fetch('/api/update-password', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id, name, note: note }),
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

  return [result, updateNote];
};

export default useUpdateNote;