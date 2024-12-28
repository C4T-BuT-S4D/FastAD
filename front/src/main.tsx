import ThemeProvider from '@mui/material/styles/ThemeProvider';
import axios from 'axios';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import { apiURL } from './config.ts';
import './index.css';
import theme from './theme.ts';

axios.defaults.baseURL = apiURL;

axios.interceptors.request.use(
  (request) => {
    return request;
  },
  (error) => {
    console.error('request error', error);
    return Promise.reject(error);
  }
);

axios.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('response error', error);
    return Promise.reject(error);
  }
);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>
  </StrictMode>,
);
