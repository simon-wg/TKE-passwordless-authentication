import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useFetchUser from "../../hooks/useFetchUser";
import NotesApp from "../../components/NotesApp";

const LoginSuccessPage = () => {
  const user = useFetchUser();

  return (
    <div>
      {/* <h1>Welcome back, {user}</h1> */}
      <NotesApp />
    </div>
  );
};

export default LoginSuccessPage;
