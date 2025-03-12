import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import config from "../config";

// This hook returns a bool whether or not the current session is authenticated
// or not.
const useAuthCheck = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const verifySession = async () => {
      try {
        const response = await fetch(
          config.backendBaseUrl + "/api/verify-session",
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
