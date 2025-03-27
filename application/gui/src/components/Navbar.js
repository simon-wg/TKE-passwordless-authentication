import React from "react";
import "./styles.css";
import useFetchUser from "../hooks/useFetchUser";
import GearIcon from "./GearIcon";
import LogoutButton from "./LogoutButton";

const Navbar = ({ setPage, currentPage }) => {
  const user = useFetchUser();

  return (
    <nav className="navbar">
      <ul className="navbar-center">
        <li>
          <button
            className={currentPage === "register" ? "active" : ""}
            onClick={() => setPage("register")}
          >
            Register
          </button>
        </li>
        <li>
          <button
            className={currentPage === "login" ? "active" : ""}
            onClick={() => setPage("login")}
          >
            Login
          </button>
        </li>
        {user !== null && (
          <>
            <li>
              <button
                className={currentPage === "app" ? "active" : ""}
                onClick={() => setPage("app")}
              >
                App
              </button>
            </li>
          </>
        )}
      </ul>

      {user !== null && (
        <div className="navbar-right">
          <LogoutButton />
          <GearIcon />
        </div>
      )}
    </nav>
  );
};

export default Navbar;
