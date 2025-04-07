import { useState } from "react";

const useCreateNote = () => {
  const [result, setResult] = useState(null);

  const createNote = async (name, note) => {
    try {
      const csrfresponse = await fetch("/api/csrf-token", {
        method: "GET",
        credentials: "include",
      });
      const token = csrfresponse.headers.get("X-CSRF-Token");

      const response = await fetch("/api/create-note", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": token,
        },
        body: JSON.stringify({ name, note }),
      });

      if (response.ok) {
        const data = await response.json();
        setResult(data);
      } else {
        setResult(false);
      }
    } catch (error) {
      console.log("Error creating note", error);
      setResult(false);
    }
  };

  return [result, createNote];
};

export default useCreateNote;
