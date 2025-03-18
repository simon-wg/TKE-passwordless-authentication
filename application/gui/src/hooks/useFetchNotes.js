import { useEffect, useState } from "react";

const useFetchNotes = () => {
  const [result, setResult] = useState([]);

    const fetchNotes = async () => {
      try {
        const response = await fetch("/api/get-user-passwords", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          console.log("Fetched from backend");
          let data = await response.json();
          data = data != null ? data : []
          setResult(result);
        }
      } catch (error) {
        console.log("Error fetching notes", error);
      }
    };

    fetchNotes();

  return result;
};

export default useFetchNotes;