import { useEffect, useState } from "react";

const useFetchNotes = () => {
  const [notes, setNotes] = useState([]);

    const fetchNotes = async () => {
      try {
        const response = await fetch("/api/get-user-passwords", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          console.log("Fetched from backend");
          const data = await response.json();
          const result = data != null ? data : []
          setNotes(result);
        }
      } catch (error) {
        console.log("Error fetching notes", error);
      }
    };

    fetchNotes();

  return notes;
};

export default useFetchNotes;