import React from 'react';
import RegisterComponent from './components/RegisterComponent';
import LoginComponent from './components/LoginComponent';

const App = () => {
  return (
    <div>
      <h1>TKey Sign in</h1>
      <RegisterComponent />
      <LoginComponent />
    </div>
  );
};

export default App;