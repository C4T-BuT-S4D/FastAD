import { centrifugeWSURL } from '@/config';
import { Service, Service_Batch } from '@/proto/data/services/services';
import { Team, Team_Batch } from '@/proto/data/teams/teams';
import { Scoreboard } from '@/proto/scoreboard/scoreboard';
import { State } from '@/proto/slac/slac';
import axios from 'axios';
import { Centrifuge, PublicationContext } from 'centrifuge';
import { useEffect, useMemo, useState } from 'react';
import { ScoreboardState } from '../models/scoreboardState';

export function useScoreboard() {
  const [teams, setTeams] = useState<Team[] | null>(null);
  const [services, setServices] = useState<Service[] | null>(null);
  const [scoreboard, setScoreboard] = useState<ScoreboardState | null>(null);

  useEffect(() => {
    async function fetchTeams() {
      const { data } = await axios.get<Team_Batch>('/teams');
      data.teams[0].name = 'very very very long team name (yes it is)';
      setTeams(data.teams);
    }

    async function fetchServices() {
      const res = await axios.get<Service_Batch>('/services');
      setServices(res.data.services);
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

      const scoreboardState = new ScoreboardState(scoreboardData);
      setScoreboard(scoreboardState);
    }

    fetchTeams();
    fetchServices();
    fetchScoreboard();
  }, []);

  useEffect(() => {
    const centrifuge = new Centrifuge(centrifugeWSURL);
    const sub = centrifuge.newSubscription('slac');
    sub.on('publication', (ctx: PublicationContext) => {
      console.log('received slac state');
      const slacState = State.fromJSON(ctx.data);

      setScoreboard((currentScoreboardState) => {
        console.log(
          'received slac state, current scoreboard state:',
          currentScoreboardState,
        );
        if (currentScoreboardState) {
          const newState = currentScoreboardState.clone();
          newState.applySlacState(slacState);
          console.log('new scoreboard state:', newState);
          return newState;
        }
        return null;
      });
    });

    // TODO: subscribe to flags.

    sub.subscribe();
    centrifuge.connect();
  }, []);

  const preparedTeams = useMemo(() => {
    if (!scoreboard || !teams) {
      return [];
    }
    return scoreboard.prepareTeams(teams);
  }, [scoreboard, teams]);

  const preparedServices = useMemo(() => {
    if (!scoreboard || !services) {
      return [];
    }
    return scoreboard.prepareServices(services);
  }, [scoreboard, services]);

  return { scoreboard, teams: preparedTeams, services: preparedServices };
}
