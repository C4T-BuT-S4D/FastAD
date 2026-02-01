import ForcADScoreboardLayout from '@/layouts/ForcADScoreboard.tsx';
import LiveLayout from '@/layouts/Live.tsx';
import TeamHistoryLayout from '@/layouts/TeamHistory.tsx';
import { setupInterceptorsTo } from '@/lib/clients/common.ts';
import theme from '@/theme.ts';
import ThemeProvider from '@mui/material/styles/ThemeProvider';
import axios from 'axios';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router';
import { apiURL } from './config.ts';
import './index.css';

axios.defaults.baseURL = apiURL;

setupInterceptorsTo(axios);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<ForcADScoreboardLayout />} />
          <Route
            path="/teams/:teamID/history"
            element={<TeamHistoryLayout />}
          />
          <Route path="/live" element={<LiveLayout />} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  </StrictMode>,
);
