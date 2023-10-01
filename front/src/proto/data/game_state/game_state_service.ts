/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import {
  GetRequest,
  GetResponse,
  UpdateRequest,
  UpdateResponse,
  UpdateRoundRequest,
  UpdateRoundResponse,
} from "./game_state";

export const protobufPackage = "data.game_state";

export type GameStateServiceDefinition = typeof GameStateServiceDefinition;
export const GameStateServiceDefinition = {
  name: "GameStateService",
  fullName: "data.game_state.GameStateService",
  methods: {
    get: {
      name: "Get",
      requestType: GetRequest,
      requestStream: false,
      responseType: GetResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          578365826: [
            new Uint8Array([17, 18, 15, 47, 97, 112, 105, 47, 103, 97, 109, 101, 95, 115, 116, 97, 116, 101]),
          ],
        },
      },
    },
    update: {
      name: "Update",
      requestType: UpdateRequest,
      requestStream: false,
      responseType: UpdateResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          578365826: [
            new Uint8Array([17, 34, 15, 47, 97, 112, 105, 47, 103, 97, 109, 101, 95, 115, 116, 97, 116, 101]),
          ],
        },
      },
    },
    updateRound: {
      name: "UpdateRound",
      requestType: UpdateRoundRequest,
      requestStream: false,
      responseType: UpdateRoundResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          578365826: [
            new Uint8Array([
              23,
              34,
              21,
              47,
              97,
              112,
              105,
              47,
              103,
              97,
              109,
              101,
              95,
              115,
              116,
              97,
              116,
              101,
              47,
              114,
              111,
              117,
              110,
              100,
            ]),
          ],
        },
      },
    },
  },
} as const;

export interface GameStateServiceImplementation<CallContextExt = {}> {
  get(request: GetRequest, context: CallContext & CallContextExt): Promise<DeepPartial<GetResponse>>;
  update(request: UpdateRequest, context: CallContext & CallContextExt): Promise<DeepPartial<UpdateResponse>>;
  updateRound(
    request: UpdateRoundRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateRoundResponse>>;
}

export interface GameStateServiceClient<CallOptionsExt = {}> {
  get(request: DeepPartial<GetRequest>, options?: CallOptions & CallOptionsExt): Promise<GetResponse>;
  update(request: DeepPartial<UpdateRequest>, options?: CallOptions & CallOptionsExt): Promise<UpdateResponse>;
  updateRound(
    request: DeepPartial<UpdateRoundRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateRoundResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | bigint | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
