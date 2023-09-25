import {createChannel, createClient} from "nice-grpc-web";
import {TeamsServiceClient, TeamsServiceDefinition} from "../../proto/data/teams/teams_service.ts";

const dataAPIAddress = 'http://localhost:11337'

const channel = createChannel(dataAPIAddress);

export const teamsService: TeamsServiceClient = createClient(TeamsServiceDefinition, channel);
