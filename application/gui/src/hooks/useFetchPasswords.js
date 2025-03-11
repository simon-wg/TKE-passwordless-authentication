import { useEffect, useState } from "react";

const useFetchPasswords = (isAuthenticated) =>  {
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
                  console.log("passwords: " + data)
                  
                }
              } catch (error) {
                console.log("Error fetching user", error);
              }
        }
        fetchPasswords()
    }, [isAuthenticated])
}

export default useFetchPasswords