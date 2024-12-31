import { useScoreboard } from '@/lib/hooks/scoreboard';
import statusColor from '@/lib/styles/statusColor';
import { Status } from '@/proto/checker/checker';
import FlagIcon from '@mui/icons-material/Flag';
import HistoryIcon from '@mui/icons-material/History';
import InfoIcon from '@mui/icons-material/Info';
import {
  Avatar,
  Box,
  IconButton,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tooltip,
  Typography,
} from '@mui/material';
import { styled, Theme } from '@mui/material/styles';
import { useWindowSize } from '@uidotdev/usehooks';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router';

const RANK_COL_WIDTH = 40;
const TEAM_COL_WIDTH = 400;
const SCORE_COL_WIDTH = 120;
const SERVICE_COL_MIN_WIDTH = 200;

const MonoFontCell = styled(TableCell)({
  fontFamily: 'Roboto Mono',
});

const HeaderCell = styled(MonoFontCell)({
  fontWeight: 500,
});

const TeamServiceCell = styled(MonoFontCell)({
  transition: 'background-color 1s ease-in-out',
  whiteSpace: 'nowrap',
  py: '0',
  px: '1em',
});

export default function ScoreboardTable() {
  const [tableWidth, setTableWidth] = useState('');
  const { width: windowWidth } = useWindowSize();
  const { teams, services, scoreboard } = useScoreboard();
  const navigate = useNavigate();

  useEffect(() => {
    const fullWidth =
      RANK_COL_WIDTH +
      TEAM_COL_WIDTH +
      SCORE_COL_WIDTH +
      SERVICE_COL_MIN_WIDTH * (services?.length ?? 1);

    if (windowWidth && windowWidth < fullWidth) {
      setTableWidth('100%');
    } else {
      setTableWidth(`${fullWidth}px`);
    }
  }, [windowWidth, services]);

  if (!teams || !services || !scoreboard || !tableWidth) {
    return <div>Loading...</div>;
  }

  return (
    <TableContainer sx={{ width: '100%', overflowX: 'auto' }}>
      <Table
        sx={{
          tableLayout: 'fixed',
          margin: '0 auto',
          width: tableWidth,
        }}
      >
        <colgroup>
          <col style={{ width: `${RANK_COL_WIDTH}px` }} />
          <col style={{ width: `${TEAM_COL_WIDTH}px` }} />
          <col style={{ width: `${SCORE_COL_WIDTH}px` }} />

          {/* Each service column divides leftover space equally */}
          {services.map((service, i) => (
            <col
              key={`service-col-${service.id}-${i}`}
              style={{ width: `${SERVICE_COL_MIN_WIDTH}px` }}
            />
          ))}
        </colgroup>

        <TableHead>
          <TableRow>
            <HeaderCell align="center"> #</HeaderCell>
            <HeaderCell align="center">Team</HeaderCell>
            <HeaderCell align="center">Score</HeaderCell>
            {services.map((s, i) => (
              <HeaderCell key={`service-header-${s.id}-${i}`} align="center">
                {s.name}
              </HeaderCell>
            ))}
          </TableRow>
        </TableHead>

        <TableBody>
          {teams.map((team, i) => {
            const rank = i + 1;
            return (
              <TableRow key={'team:' + team.id + '-' + rank}>
                <MonoFontCell align="center">{rank}</MonoFontCell>
                <MonoFontCell align="center">
                  <Stack
                    direction="row"
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 2,
                      minWidth: 0,
                    }}
                  >
                    <Avatar
                      variant="square"
                      sx={{
                        width: 64,
                        height: 64,
                      }}
                    >
                      {team.name[0].toUpperCase()}
                    </Avatar>

                    <Box
                      sx={{
                        display: 'flex',
                        flexDirection: 'column',
                        flex: 1,
                        minWidth: 0,
                      }}
                    >
                      <Typography
                        variant="subtitle1"
                        noWrap
                        sx={{
                          fontFamily: 'Roboto Mono',
                        }}
                      >
                        {team.name}
                      </Typography>

                      <Typography
                        variant="body2"
                        noWrap
                        sx={{ fontFamily: 'Roboto Mono' }}
                      >
                        {team.address}
                      </Typography>
                    </Box>
                    <Box>
                      <IconButton
                        size="small"
                        sx={{ p: 0, color: 'inherit' }}
                        onClick={() => {
                          navigate(`/teams/${team.id}/history`);
                        }}
                      >
                        <HistoryIcon />
                      </IconButton>
                    </Box>
                  </Stack>
                </MonoFontCell>
                <MonoFontCell align="center">
                  {scoreboard.getTeamScore(team.id).toFixed(2)}
                </MonoFontCell>
                {services.map((service, j) => {
                  const tss = scoreboard.getTeamServiceState(
                    team.id,
                    service.id,
                  );
                  let st: Status = Status.STATUS_UP;
                  let tooltipMessage = '';
                  if (tss?.checkStatuses && tss.checkStatuses.length > 0) {
                    st = tss.checkStatuses[tss.checkStatuses.length - 1].status;
                    tooltipMessage =
                      tss.checkStatuses[tss.checkStatuses.length - 1].message;
                  }
                  const checksTotal = parseInt(tss?.checksTotal ?? '0');
                  const checksPassed = parseInt(tss?.checksPassed ?? '0');
                  const sla =
                    (checksTotal > 0 ? checksPassed / checksTotal : 1) * 100;

                  return (
                    <TeamServiceCell
                      key={`service-${service.id}-team-${team.id}-${rank}-${j}`}
                      sx={(theme: Theme) => ({
                        backgroundColor: statusColor(theme, st).main,
                        color: '#000000',
                      })}
                    >
                      <Stack
                        direction="row"
                        gap={2}
                        alignItems="center"
                        justifyContent="space-between"
                      >
                        <Box>
                          <Typography
                            variant="body2"
                            sx={{ fontFamily: 'Roboto Mono' }}
                          >
                            <strong>SLA:</strong> {sla.toFixed(2)}%
                          </Typography>
                          <Typography
                            variant="body2"
                            sx={{ fontFamily: 'Roboto Mono' }}
                          >
                            <strong>FP:</strong> {tss?.points.toFixed(2)}
                          </Typography>
                          <Stack alignItems="center" direction="row" gap={0.5}>
                            <FlagIcon sx={{ fontSize: '1.2em' }} />
                            <Typography
                              variant="body2"
                              sx={{ fontFamily: 'Roboto Mono' }}
                            >
                              +{tss?.flagsStolen}/-{tss?.flagsLost}
                            </Typography>
                          </Stack>
                        </Box>
                        <Box sx={{ alignSelf: 'center' }}>
                          <Tooltip title={tooltipMessage}>
                            <IconButton
                              size="small"
                              sx={{
                                p: 0,
                                color: 'inherit',
                                opacity: 0.7,
                                '&:hover': {
                                  opacity: 1,
                                },
                              }}
                            >
                              <InfoIcon fontSize="small" />
                            </IconButton>
                          </Tooltip>
                        </Box>
                      </Stack>
                    </TeamServiceCell>
                  );
                })}
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
