/* eslint-disable */
import type {CallContext, CallOptions} from "nice-grpc-common";
import {SubmitFlagsRequest, SubmitFlagsResponse} from "./receiver";

export const protobufPackage = "receiver";

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

export interface ReceiverServiceImplementation<CallContextExt = {}> {
    submitFlags(
        request: SubmitFlagsRequest,
        context: CallContext & CallContextExt,
    ): Promise<DeepPartial<SubmitFlagsResponse>>;
}

export interface ReceiverServiceClient<CallOptionsExt = {}> {
    submitFlags(
        request: DeepPartial<SubmitFlagsRequest>,
        options?: CallOptions & CallOptionsExt,
    ): Promise<SubmitFlagsResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
    : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
        : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
            : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;
