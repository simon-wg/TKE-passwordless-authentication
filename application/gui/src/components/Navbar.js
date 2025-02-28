import React from 'react';
import './styles.css';

const Navbar = ({ setPage, currentPage }) => {
  return (
    <nav className="navbar">
      <ul>
        <li>
          <button
            className={currentPage === 'register' ? 'active' : ''}
            onClick={() => setPage('register')}
          >
            Register
          </button>
        </li>
        <li>
          <button
            className={currentPage === 'login' ? 'active' : ''}
            onClick={() => setPage('login')}
          >
            Login
          </button>
        </li>
      </ul>
    </nav>
  );
};

export default Navbar;