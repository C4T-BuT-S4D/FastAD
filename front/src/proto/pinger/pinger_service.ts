/* eslint-disable */
import { PingRequest, PingResponse } from "./pinger.js";

export const protobufPackage = "pinger";

export type PingerServiceDefinition = typeof PingerServiceDefinition;
export const PingerServiceDefinition = {
  name: "PingerService",
  fullName: "pinger.PingerService",
  methods: {
    ping: {
      name: "Ping",
      requestType: PingRequest,
      requestStream: false,
      responseType: PingResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;
