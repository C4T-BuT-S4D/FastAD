import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Scoreboard from './components/Scoreboard.tsx';
import Statuses from './components/Statuses.tsx';

function App() {
  return (
    <>
      <AppBar position="static">
        <h1>FastAD</h1>
      </AppBar>
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <Statuses />
      </Box>
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <Scoreboard />
      </Box>
    </>
  );
}

export default App;
