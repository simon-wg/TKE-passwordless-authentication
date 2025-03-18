// import React, { useState, useEffect } from "react";
// import RegisterComponent from "./components/RegisterComponent";
// import { Routes, Route } from "react-router-dom";
// import LoginComponent from "./components/LoginComponent";
// import Navbar from "./components/Navbar";
// import "./components/styles.css";
// import NotesApp from "./components/NotesApp";
// import useFetchUser from "./hooks/useFetchUser";

// const App = () => {
//   const [page, setPage] = useState("register");
//   const user = useFetchUser();

//   useEffect(() => {
//     if (user !== null) {
//       setPage("app");
//     }
//   }, [user]);

//   return (
//     <div>
//       <Navbar setPage={setPage} currentPage={page} />
//       {page === "register" && <RegisterComponent />}
//       {page === "login" && <LoginComponent />}
//       {page === "app" && <NotesApp />}
//     </div>
//   );
// };

// export default App;


import React, { useState, useEffect } from "react";
import RegisterComponent from "./components/RegisterComponent";
import { Routes, Route } from "react-router-dom";
import LoginComponent from "./components/LoginComponent";
import Navbar from "./components/Navbar";
import "./components/styles.css";
import NotesApp from "./components/NotesApp";
import useFetchUser from "./hooks/useFetchUser";
import LoadingCircle from "./components/LoadingCircle";

const App = () => {
  const [page, setPage] = useState("register");
  const [loading, setLoading] = useState(true);
  const user = useFetchUser();

  useEffect(() => {
    if (user !== null) {
      setPage("app");
    }
    setLoading(false);
  }, [user]);

  return (
    <div>
      <Navbar setPage={setPage} currentPage={page} />
      {loading ? (
        <LoadingCircle loading={loading} />
      ) : (
        <>
          {page === "register" && <RegisterComponent />}
          {page === "login" && <LoginComponent />}
          {page === "app" && <NotesApp />}
        </>
      )}
    </div>
  );
};

export default App;