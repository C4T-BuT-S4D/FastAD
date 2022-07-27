/* eslint-disable */
import {SubmitFlagsRequest, SubmitFlagsResponse} from "./receiver.js";

export const protobufPackage = "receiver";

/** TODO: options. */
export type ReceiverServiceDefinition = typeof ReceiverServiceDefinition;
export const ReceiverServiceDefinition = {
  name: "ReceiverService",
  fullName: "receiver.ReceiverService",
  methods: {
    submitFlags: {
      name: "SubmitFlags",
      requestType: SubmitFlagsRequest,
      requestStream: false,
      responseType: SubmitFlagsResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;
