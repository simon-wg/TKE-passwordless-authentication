import React from 'react';
import { useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCog } from '@fortawesome/free-solid-svg-icons';
import './styles.css';

const GearIcon = () => {

  const navigate = useNavigate();

  const navigateToSettings = () => {
    navigate('/settings');
  }

  return (
    <div 
      className="gear-icon"
      onClick={navigateToSettings}
    >
      <FontAwesomeIcon icon={faCog} />
    </div>
  );
};

export default GearIcon;