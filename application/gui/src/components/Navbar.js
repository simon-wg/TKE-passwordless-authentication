import React from 'react';
import './styles.css';
import useFetchUser from '../hooks/useFetchUser';

const Navbar = ({ setPage, currentPage }) => {
  const user = useFetchUser();

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
        {user !== null && (
          <li>
            <button
              className={currentPage === 'app' ? 'active' : ''}
              onClick={() => setPage('app')}
            >
              App
            </button>
          </li>
        )}
      </ul>
    </nav>
  );
};

export default Navbar;