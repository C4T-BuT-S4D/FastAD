import TeamHistoryTimeline from '@/components/ForcAD/TeamHistoryTimeline';
import Box from '@mui/material/Box';
import { Navigate, useParams } from 'react-router';

export default function TeamHistoryLayout() {
  const { teamID } = useParams();

  if (!teamID) {
    return <Navigate to="/" />;
  }

  return (
    <>
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <TeamHistoryTimeline teamID={teamID} />
      </Box>
    </>
  );
}
