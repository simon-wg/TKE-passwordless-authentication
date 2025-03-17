import React, { useState } from "react";
import RegisterComponent from "./components/RegisterComponent";
import { Routes, Route } from "react-router-dom";
import LoginComponent from "./components/LoginComponent";
import Navbar from "./components/Navbar";
import "./components/styles.css";

const App = () => {
  const [page, setPage] = useState("register");

  return (
    <div>
      <Navbar setPage={setPage} currentPage={page} />
      {page === "register" && <RegisterComponent />}
      {page === "login" && <LoginComponent />}
    </div>
  );
};

export default App;
