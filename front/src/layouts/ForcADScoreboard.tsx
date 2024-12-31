import Scoreboard from '@/components/ForcAD/Scoreboard';
import Statuses from '@/components/ForcAD/Statuses.tsx';
import TopBar from '@/components/ForcAD/TopBar';
import Box from '@mui/material/Box';

export default function ForcADScoreboardLayout() {
  return (
    <>
      <TopBar />
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <Statuses />
      </Box>
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <Scoreboard />
      </Box>
    </>
  );
}
