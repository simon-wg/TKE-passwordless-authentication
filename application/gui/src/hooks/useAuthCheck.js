import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const useAuthCheck = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const verifySession = async () => {
      try {
        const response = await fetch(
          "http://localhost:6060/api/verify_session",
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
