import { useEffect, useState } from "react";
import config from "../config";

// If the user is authenticated this hook will return the username of the
// current session user. Does this by calling the application backend.
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
