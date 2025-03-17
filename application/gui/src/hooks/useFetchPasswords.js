import { useEffect, useState } from "react";

const useFetchPasswords = (isAuthenticated) => {
  const [passwords, setPasswords] = useState([]);

  useEffect(() => {
    if (!isAuthenticated) return;

    const fetchPasswords = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/get-user-passwords", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          const data = await response.json();
          setPasswords(data);
        }
      } catch (error) {
        console.log("Error fetching passwords", error);
      }
    };

    fetchPasswords();
  }, [isAuthenticated]);

  return passwords;
};

export default useFetchPasswords;