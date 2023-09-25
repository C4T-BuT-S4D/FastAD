/* eslint-disable */
import type {CallContext, CallOptions} from "nice-grpc-common";
import {PingRequest, PingResponse} from "./pinger";

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

export interface PingerServiceImplementation<CallContextExt = {}> {
    ping(request: PingRequest, context: CallContext & CallContextExt): Promise<DeepPartial<PingResponse>>;
}

export interface PingerServiceClient<CallOptionsExt = {}> {
    ping(request: DeepPartial<PingRequest>, options?: CallOptions & CallOptionsExt): Promise<PingResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
    : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
        : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
            : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;
