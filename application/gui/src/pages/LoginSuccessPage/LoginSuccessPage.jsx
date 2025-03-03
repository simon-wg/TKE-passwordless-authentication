import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

const LoginSuccessPage = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuthentication = async () => {
      try {
        const response = await fetch("http://localhost:6060/api/verify_session", {
          method: "GET",
          credentials: "include", // Send session cookie
        });
  
        if (response.ok) {
          setIsAuthenticated(true);  // User is authenticated
        } else {
          console.log("Authentication error")
          navigate("/");  // Redirect to login page if not authenticated
        }
      } catch (error) {
        console.log("Error Authenticating")
        navigate("/");  // Redirect to login page if there's an error
      }
    };
  
    checkAuthentication();
  }, [navigate]);
  

  if (!isAuthenticated) {
    return <div>Loading...</div>;  // Or show a loading spinner
  }

  return (
    <div>
      <h1>Welcome to the Login Success Page!</h1>
      {/* Add content for the success page */}
    </div>
  );
};

export default LoginSuccessPage;
