import TopBarRoundProgress from '@/components/ForcAD/TopBarRoundProgress';
import { createProtoTransform } from '@/lib/clients/common';
import { GameState } from '@/proto/data/game_state/game_state';
import { AppBar, Box, Toolbar, Typography } from '@mui/material';
import axios from 'axios';
import { useEffect, useState } from 'react';
import { Link } from 'react-router';

export default function TopBar() {
  const [gameState, setGameState] = useState<GameState | null>(null);

  useEffect(() => {
    async function fetchGameState() {
      const { data } = await axios.get<GameState>('/game', {
        transformResponse: createProtoTransform(GameState),
        responseType: 'arraybuffer',
        params: {
          format: 'proto',
        },
      });
      setGameState(data);
    }

    fetchGameState();
    const interval = setInterval(fetchGameState, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <AppBar position="static" color="default">
      <Toolbar>
        <Typography
          variant="h6"
          component="div"
          sx={{ flexGrow: 0, fontFamily: 'Roboto Mono', mr: 4 }}
        >
          FastAD
        </Typography>

        <Box
          component={Link}
          to="/live"
          sx={{ textDecoration: 'none', color: 'inherit', mr: 2 }}
        >
          <Typography variant="body1" sx={{ fontFamily: 'Roboto Mono' }}>
            Live
          </Typography>
        </Box>

        <Box
          component={Link}
          to="https://github.com/C4T-BuT-S4D/FastAD"
          sx={{ textDecoration: 'none', color: 'inherit' }}
        >
          <Typography variant="body1" sx={{ fontFamily: 'Roboto Mono' }}>
            GitHub
          </Typography>
        </Box>

        <Box sx={{ flexGrow: 1 }} />

        {gameState?.runningRoundStart && gameState?.roundDuration && (
          <TopBarRoundProgress
            round={gameState.runningRound}
            totalRounds={gameState.totalRounds}
            startTime={gameState.runningRoundStart}
            duration={gameState.roundDuration}
          />
        )}
      </Toolbar>
    </AppBar>
  );
}
