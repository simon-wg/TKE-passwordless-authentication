import { useEffect, useState } from "react";
import { secureFetch } from "../util/secureFetch";

/**
 * Custom hook to fetch user notes from the server.
 *
 * @returns {Array} An array of user notes.
 *
 * @example
 * const notes = useFetchNotes();
 * console.log(notes);
 *
 * @remarks
 * This hook fetches notes from the endpoint "/api/get-user-note" using a GET request.
 * It includes credentials in the request and handles the response by setting the result state.
 * If the response is not ok or an error occurs, it logs the error to the console.
 */
const useFetchNotes = () => {
  const [result, setResult] = useState([]);

  useEffect(() => {
    const fetchNotes = async () => {
      try {
        const response = await secureFetch("/api/get-user-note");

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
