import { useState } from "react";
import { secureFetch } from "../util/secureFetch";
/**
 * Custom hook to delete a note.
 *
 * @returns {[any, Function]} An array containing the delete result and the deleteNote function.
 *
 * @example
 * const [deleteResult, deleteNote] = useDeleteNote();
 *
 * // To delete a note
 * deleteNote(noteId);
 *
 * @typedef {Object} DeleteResult
 * @property {boolean} success - Indicates if the deletion was successful.
 * @property {string} message - A message related to the deletion result.
 *
 * @function deleteNote
 * @param {string} id - The ID of the note to be deleted.
 * @returns {Promise<void>} A promise that resolves when the note is deleted.
 */
const useDeleteNote = () => {
  const [deleteResult, setDeleteResult] = useState(null);

  const deleteNote = async (id) => {
    try {
      const response = await secureFetch("/api/delete-note", {
        method: "DELETE",
        body: JSON.stringify({ id }),
      });

      if (response.ok) {
        const data = await response.json();
        setDeleteResult(data);
      } else {
        setDeleteResult(false);
      }
    } catch (error) {
      console.log("Error deleting note", error);
      setDeleteResult(false);
    }
  };

  return [deleteResult, deleteNote];
};

export default useDeleteNote;
