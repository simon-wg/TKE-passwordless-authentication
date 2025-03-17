import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useFetchUser from "../../hooks/useFetchUser";

const LoginSuccessPage = () => {
  const user = useFetchUser();
  const navigate = useNavigate();

  if (!user) {
    navigate("/");
  }

  return (
    <div>
      <h1>Welcome back, {user}</h1>
    </div>
  );
};

export default LoginSuccessPage;
