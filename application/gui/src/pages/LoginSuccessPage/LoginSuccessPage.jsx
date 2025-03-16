import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useAuthCheck from "../../hooks/useAuthCheck";
import useFetchUser from "../../hooks/useFetchUser";

const LoginSuccessPage = () => {
  const user = useFetchUser();

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Welcome back, {user}</h1>
    </div>
  );
};

export default LoginSuccessPage;
