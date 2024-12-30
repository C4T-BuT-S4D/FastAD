import Scoreboard from '@/components/ForcAD/Scoreboard.tsx';
import Statuses from '@/components/ForcAD/Statuses.tsx';
import Box from '@mui/material/Box';

export default function ForcADLayout() {
  return (
    <>
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <Statuses />
      </Box>
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <Scoreboard />
      </Box>
    </>
  );
}
