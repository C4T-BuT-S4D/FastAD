import ThemeProvider from '@mui/material/styles/ThemeProvider';
import axios from 'axios';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import { apiURL } from './config.ts';
import './index.css';
import { setupInterceptorsTo } from './lib/clients/common.ts';
import theme from './theme.ts';

axios.defaults.baseURL = apiURL;

setupInterceptorsTo(axios);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>
  </StrictMode>,
);
