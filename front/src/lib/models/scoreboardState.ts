import { Service } from '@/proto/data/services/services';
import { Team } from '@/proto/data/teams/teams';
import { Scoreboard } from '@/proto/scoreboard/scoreboard';

import { Scoreboard_TeamServiceState } from '@/proto/scoreboard/scoreboard';
import { State as SlacState, TeamServiceState } from '@/proto/slac/slac';

export class ScoreboardState {
  private teamServiceState: Map<
    string,
    Map<string, Scoreboard_TeamServiceState>
  >;
  private teamScore: Map<string, number>;

  constructor(scoreboard: Scoreboard) {
    this.teamServiceState = new Map();
    this.teamScore = new Map();

    scoreboard.teamServiceStates.forEach((tss) => {
      if (!this.teamServiceState.has(tss.teamId)) {
        const serviceState = new Map<string, Scoreboard_TeamServiceState>();
        serviceState.set(tss.serviceId, tss);
        this.teamServiceState.set(tss.teamId, serviceState);
      } else {
        this.teamServiceState.get(tss.teamId)?.set(tss.serviceId, tss);
      }
    });
    this.calculateScores();
  }

  public clone() {
    const newState = new ScoreboardState({
      teamServiceStates: [],
    });

    // Deep clone the team service state map
    this.teamServiceState.forEach((serviceMap, teamId) => {
      const newServiceMap = new Map<string, Scoreboard_TeamServiceState>();
      serviceMap.forEach((tss, serviceId) => {
        newServiceMap.set(serviceId, {
          ...tss,
          checkStatuses: [...tss.checkStatuses],
        });
      });
      newState.teamServiceState.set(teamId, newServiceMap);
    });

    // Clone the team score map
    this.teamScore.forEach((score, teamId) => {
      newState.teamScore.set(teamId, score);
    });

    return newState;
  }

  public applySlacState(slacState: SlacState) {
    function newServiceState(tss: TeamServiceState) {
      return {
        teamId: tss.teamId,
        serviceId: tss.serviceId,
        checksTotal: tss.checksTotal,
        checksPassed: tss.checksPassed,
        checkStatuses: tss.checkStatuses,
        points: 0,
        flagsStolen: '0',
        flagsLost: '0',
      };
    }

    slacState.teamServiceStates.forEach((tss) => {
      const currentTeamState = this.teamServiceState.get(tss.teamId);
      const currentServiceState = currentTeamState?.get(tss.serviceId);

      if (currentTeamState && currentServiceState) {
        currentTeamState.set(tss.serviceId, {
          ...currentServiceState,
          checksPassed: tss.checksPassed,
          checksTotal: tss.checksTotal,
          checkStatuses: tss.checkStatuses,
        });
      } else if (currentTeamState) {
        currentTeamState.set(tss.serviceId, newServiceState(tss));
      } else {
        const serviceState = new Map<string, Scoreboard_TeamServiceState>();
        serviceState.set(tss.serviceId, newServiceState(tss));
        this.teamServiceState.set(tss.teamId, serviceState);
      }
    });
    this.calculateScores();
  }

  private calculateScores() {
    this.teamScore.clear();
    this.teamServiceState.forEach((serviceState) => {
      serviceState.forEach((tss) => {
        const checksTotal = parseInt(tss.checksTotal);
        const checksPassed = parseInt(tss.checksPassed);
        this.teamScore.set(
          tss.teamId,
          (this.teamScore.get(tss.teamId) ?? 0) +
            tss.points * (checksTotal > 0 ? checksPassed / checksTotal : 1),
        );
      });
    });
  }

  public prepareTeams(teams: Team[]) {
    const sortedTeams = teams
      .map((t) => ({
        ...t,
      }))
      .sort((a, b) => {
        const scoreA = this.teamScore.get(a.id) ?? 0;
        const scoreB = this.teamScore.get(b.id) ?? 0;
        return scoreB - scoreA;
      });
    return Array(50).fill(sortedTeams).flat();
  }

  public prepareServices(services: Service[]) {
    const sortedServices = [...services].sort((a, b) => {
      return parseInt(a.id) - parseInt(b.id);
    });
    return Array(1).fill(sortedServices).flat();
  }

  public getTeamScore(teamId: string) {
    return this.teamScore.get(teamId) ?? 0;
  }

  public getTeamServiceState(teamId: string, serviceId: string) {
    return this.teamServiceState.get(teamId)?.get(serviceId);
  }
}
