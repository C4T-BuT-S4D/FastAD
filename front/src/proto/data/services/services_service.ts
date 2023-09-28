/* eslint-disable */
import type {CallContext, CallOptions} from "nice-grpc-common";
import {ListRequest, ListResponse} from "./services";

export const protobufPackage = "data.services";

export type ServicesServiceDefinition = typeof ServicesServiceDefinition;
export const ServicesServiceDefinition = {
    name: "ServicesService",
    fullName: "data.services.ServicesService",
    methods: {
        list: {
            name: "List",
            requestType: ListRequest,
            requestStream: false,
            responseType: ListResponse,
            responseStream: false,
            options: {},
        },
    },
} as const;

export interface ServicesServiceImplementation<CallContextExt = {}> {
    list(request: ListRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ListResponse>>;
}

export interface ServicesServiceClient<CallOptionsExt = {}> {
    list(request: DeepPartial<ListRequest>, options?: CallOptions & CallOptionsExt): Promise<ListResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | bigint | undefined;

export type DeepPartial<T> = T extends Builtin ? T
    : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
        : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
            : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;
