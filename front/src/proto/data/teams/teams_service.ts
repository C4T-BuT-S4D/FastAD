/* eslint-disable */
import type {CallContext, CallOptions} from "nice-grpc-common";
import {CreateBatchRequest, CreateBatchResponse, ListRequest, ListResponse} from "./teams";

export const protobufPackage = "data.teams";

export type TeamsServiceDefinition = typeof TeamsServiceDefinition;
export const TeamsServiceDefinition = {
    name: "TeamsService",
    fullName: "data.teams.TeamsService",
    methods: {
        list: {
            name: "List",
            requestType: ListRequest,
            requestStream: false,
            responseType: ListResponse,
            responseStream: false,
            options: {},
        },
        createBatch: {
            name: "CreateBatch",
            requestType: CreateBatchRequest,
            requestStream: false,
            responseType: CreateBatchResponse,
            responseStream: false,
            options: {},
        },
    },
} as const;

export interface TeamsServiceImplementation<CallContextExt = {}> {
    list(request: ListRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ListResponse>>;

    createBatch(
        request: CreateBatchRequest,
        context: CallContext & CallContextExt,
    ): Promise<DeepPartial<CreateBatchResponse>>;
}

export interface TeamsServiceClient<CallOptionsExt = {}> {
    list(request: DeepPartial<ListRequest>, options?: CallOptions & CallOptionsExt): Promise<ListResponse>;

    createBatch(
        request: DeepPartial<CreateBatchRequest>,
        options?: CallOptions & CallOptionsExt,
    ): Promise<CreateBatchResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
    : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
        : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
            : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;
