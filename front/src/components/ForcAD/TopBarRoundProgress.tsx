import { Duration } from '@/proto/google/protobuf/duration';
import { Box, CircularProgress, Stack, Typography } from '@mui/material';

import { useEffect, useState } from 'react';

interface Props {
  round: string;
  totalRounds: string | undefined;
  startTime: Date;
  duration: Duration;
}

export default function TopBarRoundProgress(props: Props) {
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    const updateProgress = () => {
      const nowNanos = Date.now() * 1e6;
      const startTimeNanos = props.startTime.getTime() * 1e6;
      const durationNanos =
        parseInt(props.duration.seconds) * 1e9 + props.duration.nanos;
      const elapsedNanos = nowNanos - startTimeNanos;
      const percentage = (elapsedNanos / durationNanos) * 100;
      setProgress(Math.min(percentage, 100));
    };

    const timer = setInterval(updateProgress, 100);
    updateProgress(); // Initial update

    return () => clearInterval(timer);
  }, [props]);

  return (
    <Stack direction="row" spacing={1} alignItems="center">
      <Typography variant="caption" component="div" color="text.secondary">
        {props.totalRounds && props.totalRounds !== '0'
          ? `Round ${props.round}/${props.totalRounds}`
          : `Round ${props.round}`}
      </Typography>
      <Box sx={{ position: 'relative', display: 'inline-flex' }}>
        <CircularProgress variant="determinate" value={progress} />
        <Box
          sx={{
            top: 0,
            left: 0,
            bottom: 0,
            right: 0,
            position: 'absolute',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Typography variant="caption" component="div" color="text.secondary">
            {`${Math.round(progress)}%`}
          </Typography>
        </Box>
      </Box>
    </Stack>
  );
}
