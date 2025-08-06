import React from 'react'; 
import ReactDOM from 'react-dom/client';
import App from './App.jsx';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')).render(
  // Now the file knows what "React" is, and this line will work
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)