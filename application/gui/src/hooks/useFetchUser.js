import { useEffect, useState } from "react";

const useFetchUser = (isAuthenticated) => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    if (!isAuthenticated) return;

    const fetchUser = async () => {
      try {
        const response = await fetch("http://localhost:6060/getuser", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          const data = await response.json();
          setUser(data.user);
        }
      } catch (error) {
        console.log("Error fetching user", error);
      }
    };

    fetchUser();
  }, [isAuthenticated]);

  return user;
};

export default useFetchUser;
