import { useState } from 'react';

const useDeleteNote = () => {
  const [result, setResult] = useState(null);

  const deleteNote = async (id) => {
    try {
      const response = await fetch('/api/delete-note', {
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
      console.log('Error deleting note', error);
      setResult(false);
    }
  };

  return [result, deleteNote];
};

export default useDeleteNote;