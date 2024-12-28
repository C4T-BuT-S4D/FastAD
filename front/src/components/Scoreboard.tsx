import FlagIcon from '@mui/icons-material/Flag';
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
  Typography
} from '@mui/material';
import { Theme } from '@mui/material/styles';
import { useMeasure } from '@uidotdev/usehooks';
import axios from 'axios';
import { useEffect, useState } from 'react';
import { Status } from '../proto/checker/checker';
import { Service, Service_Batch } from '../proto/data/services/services';
import { Team, Team_Batch } from '../proto/data/teams/teams';
import { Scoreboard, Scoreboard_TeamServiceState } from '../proto/scoreboard/scoreboard';

const RANK_COL_WIDTH = 50;
const TEAM_COL_WIDTH = 400;
const SCORE_COL_WIDTH = 150;
const SERVICE_COL_MIN_WIDTH = 230;

const statusColor = (theme: Theme, status: Status) => {
  // Convert the numeric `status` enum to the typed Status enum
  switch (status) {
    case Status.STATUS_UP:
      return theme.palette.statusUp?.main;
    case Status.STATUS_DOWN:
      return theme.palette.statusDown?.main;
    case Status.STATUS_CORRUPT:
      return theme.palette.statusCorrupt?.main;
    case Status.STATUS_MUMBLE:
      return theme.palette.statusMumble?.main;
    case Status.STATUS_CHECK_FAILED:
      return theme.palette.statusCheckFailed?.main;
    default:
      return theme.palette.grey[200]; // Fallback background
  }
};

export default function ScoreboardTable() {
  const [teams, setTeams] = useState<Team[]>([]);
  const [teamsLoaded, setTeamsLoaded] = useState(false);

  const [services, setServices] = useState<Service[]>([]);
  const [servicesLoaded, setServicesLoaded] = useState(false);

  const [teamServiceState, setTeamServiceState] = useState<
    Record<string, Record<string, Scoreboard_TeamServiceState>>
  >({});
  const [teamScore, setTeamScore] = useState<Record<string, number>>({});
  const [scoreboardLoaded, setScoreboardLoaded] = useState(false);

  const [serviceColWidth, setServiceColWidth] = useState('200px');
  const [tableRef, { width: tableWidth }] = useMeasure();

  useEffect(() => {
    async function fetchTeams() {
      const { data } = await axios.get<Team_Batch>('/teams');
      data.teams[0].name = 'very very very long team name (yes it is)';
      setTeams(data.teams);
      setTeamsLoaded(true);
    }

    async function fetchServices() {
      const res = await axios.get<Service_Batch>('/services');
      setServices(res.data.services);
      setServicesLoaded(true);
    }

    async function fetchScoreboard() {
      const { data: scoreboardData } = await axios.get<Scoreboard>(
        '/scoreboard',
        {
          transformResponse: (resp) => {
            return Scoreboard.fromJSON(JSON.parse(resp));
          },
        },
      );

      const teamServiceStatus: Record<
        string,
        Record<string, Scoreboard_TeamServiceState>
      > = {};
      const teamScoreTemp: Record<string, number> = {};

      scoreboardData.teamServiceStates.forEach((tss) => {
        if (!teamServiceStatus[tss.teamId]) {
          teamServiceStatus[tss.teamId] = {};
        }
        if (!teamScoreTemp[tss.teamId]) {
          teamScoreTemp[tss.teamId] = 0;
        }

        teamServiceStatus[tss.teamId][tss.serviceId] = tss;

        const checksTotal = parseInt(tss.checksTotal);
        const checksPassed = parseInt(tss.checksPassed);
        teamScoreTemp[tss.teamId] +=
          tss.points * (checksTotal > 0 ? checksPassed / checksTotal : 1);
      });

      setTeamServiceState(teamServiceStatus);
      setTeamScore(teamScoreTemp);
      setScoreboardLoaded(true);
    }

    fetchTeams();
    fetchServices();
    fetchScoreboard();
  }, []);

  useEffect(() => {
    const leftoverPx = RANK_COL_WIDTH + TEAM_COL_WIDTH + SCORE_COL_WIDTH;
    // Calculate minimum width needed for all columns
    const serviceWidth =
      services.length > 0 && tableWidth
        ? Math.min(
          SERVICE_COL_MIN_WIDTH,
          (tableWidth - leftoverPx) / services.length
        )
        : SERVICE_COL_MIN_WIDTH;

    console.log('tableWidth', tableWidth);
    console.log('services.length', services.length);
    console.log('leftoverPx', leftoverPx);
    console.log('serviceWidth', serviceWidth);

    setServiceColWidth(`${serviceWidth}px`);
  }, [services, tableWidth]);

  if (!teamsLoaded || !servicesLoaded || !scoreboardLoaded) {
    return <div>Loading...</div>;
  }

  // Optionally, sort teams by highest score first
  let sortedTeams = [...teams].sort((a, b) => {
    const scoreA = teamScore[a.id] ?? 0;
    const scoreB = teamScore[b.id] ?? 0;
    return scoreB - scoreA;
  });

  // Create an array of 1000 copies of the sorted teams
  const expandedTeams = Array(50).fill(sortedTeams).flat();
  sortedTeams = expandedTeams;

  let sortedServices = [...services].sort((a, b) => {
    return parseInt(a.id) - parseInt(b.id);
  });
  sortedServices = Array(4).fill(sortedServices).flat();

  return (
    <TableContainer sx={{ maxWidth: '100%', overflowX: 'auto' }} ref={tableRef}>
      <Table
        sx={{
          tableLayout: 'fixed',
          margin: '0 auto',
          maxWidth: '100%',
          width: 'auto'
        }}
      >
        <colgroup>
          <col style={{ width: `${RANK_COL_WIDTH}px` }} />
          <col style={{ width: `${TEAM_COL_WIDTH}px` }} />
          <col style={{ width: `${SCORE_COL_WIDTH}px` }} />

          {/* Each service column divides leftover space equally */}
          {sortedServices.map((service) => (
            <col
              key={'col-service-' + service.id}
              style={{ width: `${SERVICE_COL_MIN_WIDTH}px` }}
            />
          ))}
        </colgroup>

        <TableHead>
          <TableRow>
            <TableCell align="center" sx={{ fontFamily: 'Roboto Mono' }}>
              #
            </TableCell>
            <TableCell align="center" sx={{ fontFamily: 'Roboto Mono' }}>
              Team
            </TableCell>
            <TableCell align="center" sx={{ fontFamily: 'Roboto Mono' }}>
              Score
            </TableCell>
            {sortedServices.map((s) => (
              <TableCell
                key={'hdr-' + s.id}
                align="center"
                sx={{ fontFamily: 'Roboto Mono' }}
              >
                {s.name}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>

        <TableBody>
          {sortedTeams.map((team, i) => {
            const rank = i + 1;
            return (
              <TableRow key={'team:' + team.id + '-' + rank}>
                <TableCell align="center" sx={{ fontFamily: 'Roboto Mono' }}>
                  {rank}
                </TableCell>
                <TableCell align="center" sx={{ fontFamily: 'Roboto Mono' }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Avatar
                      variant="square"
                      sx={{
                        width: '64px',
                        height: '64px',
                        fontFamily: 'Roboto Mono'
                      }}
                    >
                      {team.name[0]}
                    </Avatar>
                    <Box sx={{ flex: 1 }}>
                      <Typography variant="subtitle1">{team.name}</Typography>
                      <Typography
                        variant="body2"
                        sx={{ fontFamily: 'Roboto Mono' }}
                      >
                        {team.address}
                      </Typography>
                    </Box>
                  </Box>
                </TableCell>
                <TableCell align="center" sx={{ fontFamily: 'Roboto Mono' }}>
                  {(teamScore[team.id] ?? 0).toFixed(2)}
                </TableCell>
                {sortedServices.map((service) => {
                  const tss = teamServiceState[team.id]?.[service.id];
                  const st = tss?.status ?? Status.STATUS_UNSPECIFIED;
                  const bg = (theme: Theme) => statusColor(theme, st);
                  const checksTotal = parseInt(tss?.checksTotal ?? '0');
                  const checksPassed = parseInt(tss?.checksPassed ?? '0');
                  const sla =
                    (checksTotal > 0 ? checksPassed / checksTotal : 1) * 100;

                  return (
                    <TableCell
                      key={`service-${service.id}-team-${team.id}`}
                      sx={(theme: Theme) => ({
                        backgroundColor: bg(theme),
                        color: theme.palette.getContrastText(bg(theme) || ''),
                        whiteSpace: 'nowrap',
                        paddingLeft: '1em',
                        py: '0'
                      })}
                    >
                      <Stack
                        direction="row"
                        gap={0.5}
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
                          <Tooltip title="Service Details">
                            <IconButton
                              size="small"
                              onClick={() => {
                              }}
                              sx={{
                                color: 'inherit',
                                opacity: 0.7,
                                '&:hover': {
                                  opacity: 1
                                },
                              }}
                            >
                              <InfoIcon fontSize="small" />
                            </IconButton>
                          </Tooltip>
                        </Box>
                      </Stack>
                    </TableCell>
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
