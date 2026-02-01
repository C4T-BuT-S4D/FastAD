import TeamHistoryTimeline from '@/components/ForcAD/TeamHistoryTimeline';
import TopBar from '@/components/ForcAD/TopBar';
import Box from '@mui/material/Box';
import { Navigate, useParams, useSearchParams } from 'react-router';

export default function TeamHistoryLayout() {
  const { teamID } = useParams();
  const [searchParams] = useSearchParams();
  const serviceID = searchParams.get('serviceId') ?? '';

  if (!teamID) {
    return <Navigate to="/" />;
  }

  return (
    <>
      <TopBar />
      <Box sx={{ paddingTop: 2, paddingLeft: 2, paddingRight: 2 }}>
        <TeamHistoryTimeline teamID={teamID} serviceID={serviceID} />
      </Box>
    </>
  );
}
