import { useEffect, useState } from "react";
import config from "../config";

/**
 * Custom hook to fetch the current user data from the server.
 *
 * This hook sends a GET request to the "/api/getuser" endpoint to retrieve
 * the user data. The request includes credentials (cookies) for authentication.
 * If the request is successful, the user data is stored in the state.
 *
 * @returns {Object|null} The user data if available, otherwise null.
 */
const useFetchUser = () => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await fetch("/api/getuser", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          const data = await response.json();
          setUser(data.user);
        }
      } catch (error) {
        console.log("Error fetching user");
      }
    };

    fetchUser();
  }, []);

  return user;
};

export default useFetchUser;
