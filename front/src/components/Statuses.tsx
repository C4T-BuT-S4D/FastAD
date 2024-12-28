import Stack from '@mui/material/Stack';
import { styled } from '@mui/material/styles';

const StatusBarItem = styled('div')({
  width: '100%',
  textAlign: 'center',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'black',
  padding: '0.3em'
});

const UpStatusBarItem = styled(StatusBarItem)(({ theme }) => ({
  backgroundColor: theme.palette.statusUp!.main
}));

const CorruptStatusBarItem = styled(StatusBarItem)(({ theme }) => ({
  backgroundColor: theme.palette.statusCorrupt!.main
}));

const MumbleStatusBarItem = styled(StatusBarItem)(({ theme }) => ({
  backgroundColor: theme.palette.statusMumble!.main
}));

const DownStatusBarItem = styled(StatusBarItem)(({ theme }) => ({
  backgroundColor: theme.palette.statusDown!.main
}));

const CheckFailedStatusBarItem = styled(StatusBarItem)(({ theme }) => ({
  backgroundColor: theme.palette.statusCheckFailed!.main
}));

export default function Statuses() {
  return (
    <Stack direction={'row'} spacing={0}>
      <UpStatusBarItem>UP</UpStatusBarItem>
      <MumbleStatusBarItem>MUMBLE</MumbleStatusBarItem>
      <CorruptStatusBarItem>CORRUPT</CorruptStatusBarItem>
      <DownStatusBarItem>DOWN</DownStatusBarItem>
      <CheckFailedStatusBarItem>CHECK FAILED</CheckFailedStatusBarItem>
    </Stack>
  );
}
