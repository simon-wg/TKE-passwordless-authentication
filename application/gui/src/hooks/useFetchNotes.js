import { useEffect, useState } from "react";

const useFetchNotes = (isAuthenticated) => {
  const [notes, setNotes] = useState([]);

  useEffect(() => {
    if (!isAuthenticated) return;

    const fetchNotes = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/get-user-passwords", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          const data = await response.json();
          setNotes(data);
        }
      } catch (error) {
        console.log("Error fetching notes", error);
      }
    };

    fetchNotes();
  }, [isAuthenticated]);

  return notes;
};

export default useFetchNotes;