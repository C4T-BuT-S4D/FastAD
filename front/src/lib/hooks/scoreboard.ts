import { centrifugeWSURL } from '@/config';
import { Service, Service_Batch } from '@/proto/data/services/services';
import { Team, Team_Batch } from '@/proto/data/teams/teams';
import { Scoreboard } from '@/proto/scoreboard/scoreboard';
import axios from 'axios';
import { Centrifuge, PublicationContext } from 'centrifuge';
import { useEffect, useMemo, useState } from 'react';
import { createProtoTransform } from '../clients/common';
import { ScoreboardState } from '../models/scoreboardState';

export function useScoreboard() {
  const [teams, setTeams] = useState<Team[] | null>(null);
  const [services, setServices] = useState<Service[] | null>(null);
  const [scoreboard, setScoreboard] = useState<ScoreboardState | null>(null);

  useEffect(() => {
    async function fetchTeams() {
      const { data } = await axios.get<Team_Batch>('/teams', {
        transformResponse: createProtoTransform(Team_Batch),
      });
      data.teams[0].name = 'very very very long team name (yes it is)';
      setTeams(data.teams);
    }

    async function fetchServices() {
      const res = await axios.get<Service_Batch>('/services', {
        transformResponse: createProtoTransform(Service_Batch),
      });
      setServices(res.data.services);
    }

    async function fetchScoreboard() {
      const { data: scoreboardData } = await axios.get<Scoreboard>(
        '/scoreboard',
        {
          transformResponse: createProtoTransform(Scoreboard),
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

    const scoreboardSub = centrifuge.newSubscription('scoreboard');
    scoreboardSub.on('publication', (ctx: PublicationContext) => {
      const scoreboardState = new ScoreboardState(
        Scoreboard.fromJSON(ctx.data),
      );
      setScoreboard(scoreboardState);
    });

    scoreboardSub.subscribe();
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
