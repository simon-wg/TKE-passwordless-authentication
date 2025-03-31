import React, { useState, useEffect } from "react";
import RegisterComponent from "./components/RegisterComponent";
import { Routes, Route } from "react-router-dom";
import LoginComponent from "./components/LoginComponent";
import Navbar from "./components/Navbar";
import "./components/styles.css";
import NotesApp from "./components/NotesApp";
import useFetchUser from "./hooks/useFetchUser";
import LoadingCircle from "./components/LoadingCircle";
import UnAuth from "./components/UnAuth";

const App = () => {
  const [page, setPage] = useState("register");
  const [loading, setLoading] = useState(true);
  const user = useFetchUser();

  useEffect(() => {
    if (user !== null) {
      setPage("app");
    }else {
      setPage("start");
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
          {page === "start" && <UnAuth setPage={setPage} />}
        </>
      )}
    </div>
  );
};

export default App;