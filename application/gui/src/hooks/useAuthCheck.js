import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

/**
 * A custom hook that uses the api /api/verify-session to determine if the user has a valid session cookie
 * If a user is not authenticated or an error occurs they are redirected to "/"
 * @returns bool
 *
 */
const useAuthCheck = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const verifySession = async () => {
      try {
        const response = await fetch(
          "http://localhost:8080/api/verify-session",
          {
            method: "GET",
            credentials: "include",
          }
        );

        if (!response.ok) {
          console.log("Authentication error");
          navigate("/");
          return;
        }

        setIsAuthenticated(true);
      } catch (error) {
        console.log("Error Authenticating", error);
        navigate("/");
      }
    };

    verifySession();
  }, [navigate]);

  return isAuthenticated;
};

export default useAuthCheck;
