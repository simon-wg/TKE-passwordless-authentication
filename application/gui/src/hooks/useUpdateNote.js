import { useState } from "react";
import { secureFetch } from "../util/secureFetch";

/**
 * Custom hook to update a note.
 *
 * @returns {[boolean|null, Function]} - Returns an array with the result of the update operation and the updateNote function.
 *
 * @example
 * const [result, updateNote] = useUpdateNote();
 *
 * // To update a note
 * updateNote(id, name, note);
 *
 * @function
 * @name useUpdateNote
 *
 * @async
 * @param {string} id - The ID of the note to update.
 * @param {string} name - The name of the note.
 * @param {string} note - The content of the note.
 */
const useUpdateNote = () => {
  const [result, setResult] = useState(null);

  const updateNote = async (id, name, note) => {
    try {
      const response = await secureFetch("/api/update-note", {
        method: "POST",
        body: JSON.stringify({ id, name, note: note }),
      });

      if (response.ok) {
        setResult(true);
      } else {
        setResult(false);
      }
    } catch (error) {
      console.log("Error saving note", error);
      setResult(false);
    }
  };

  return [result, updateNote];
};

export default useUpdateNote;
