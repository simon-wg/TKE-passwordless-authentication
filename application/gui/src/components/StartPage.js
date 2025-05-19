import React from "react";
import "./StartPage.css";

const StartPage = ({ setPage }) => {
  const handleRegisterClick = () => {
    setPage("register");
  };

  return (
    <div className="centered-container">
      {" "}
      <div className="get-started-box">
        {" "}
        <h1 style={{ textAlign: "center" }}>
          Welcome to the TKey Login Demo
        </h1>{" "}
        <div className="get-started-container">
          {" "}
          <h2>Getting Started</h2>{" "}
          <ul>
            {" "}
            <li>
              {" "}
              Make sure you have the{" "}
              <a href="https://github.com/epicreach/tkey-web-authentication">
                daemon
              </a>{" "}
              installed and running{" "}
            </li>{" "}
            <li>Go to the Register page and fill out the forms</li>{" "}
            <li>
              {" "}
              When registering or logging in, you will be prompted to enter an
              optional USS (User Supplied Secret){" "}
            </li>{" "}
            <li>
              {" "}
              If you've created a USS you will need to enter the same USS each
              time you login{" "}
            </li>{" "}
          </ul>{" "}
        </div>{" "}
        <button className="to-register-btn" onClick={handleRegisterClick}>
          {" "}
          Register with TKey{" "}
        </button>{" "}
      </div>{" "}
    </div>
  );
};

export default StartPage;
