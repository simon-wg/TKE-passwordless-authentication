import { useEffect, useState } from "react";

const useFetchNotes = () => {
  const [result, setResult] = useState([]);

  useEffect(() => {
    const fetchNotes = async () => {
      try {
        const response = await fetch("/api/get-user-passwords", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          let data = await response.json();
          data = data != null ? data : [];
          setResult(data);
        }
      } catch (error) {
        console.log("Error fetching notes", error);
      }
    };

    fetchNotes();
  }, []);

  return result;
};

export default useFetchNotes;