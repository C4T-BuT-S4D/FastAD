/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Duration } from "../../google/protobuf/duration";
import { Timestamp } from "../../google/protobuf/timestamp";
import { Version } from "../version/version";

export const protobufPackage = "data.game_state";

export const GameMode = { GAME_MODE_UNSPECIFIED: 0, GAME_MODE_CLASSIC: 1 } as const;

export type GameMode = typeof GameMode[keyof typeof GameMode];

export function gameModeFromJSON(object: any): GameMode {
  switch (object) {
    case 0:
    case "GAME_MODE_UNSPECIFIED":
      return GameMode.GAME_MODE_UNSPECIFIED;
    case 1:
    case "GAME_MODE_CLASSIC":
      return GameMode.GAME_MODE_CLASSIC;
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum GameMode");
  }
}

export function gameModeToJSON(object: GameMode): string {
  switch (object) {
    case GameMode.GAME_MODE_UNSPECIFIED:
      return "GAME_MODE_UNSPECIFIED";
    case GameMode.GAME_MODE_CLASSIC:
      return "GAME_MODE_CLASSIC";
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum GameMode");
  }
}

export interface GameState {
  startTime: Date | undefined;
  endTime: Date | undefined;
  totalRounds: number;
  paused: boolean;
  flagLifetimeRounds: number;
  roundDuration: Duration | undefined;
  mode: GameMode;
  runningRound: number;
  runningRoundStart: Date | undefined;
}

export interface GetRequest {
  version: Version | undefined;
}

export interface GetResponse {
  gameState: GameState | undefined;
  version: Version | undefined;
}

export interface UpdateRequest {
  startTime: Date | undefined;
  endTime: Date | undefined;
  totalRounds: number;
  paused: boolean;
  flagLifetimeRounds: number;
  roundDuration: Duration | undefined;
  mode: GameMode;
}

export interface UpdateResponse {
  gameState: GameState | undefined;
  version: Version | undefined;
}

export interface UpdateRoundRequest {
  runningRound: number;
  runningRoundStart: Date | undefined;
}

export interface UpdateRoundResponse {
  gameState: GameState | undefined;
  version: Version | undefined;
}

function createBaseGameState(): GameState {
  return {
    startTime: undefined,
    endTime: undefined,
    totalRounds: 0,
    paused: false,
    flagLifetimeRounds: 0,
    roundDuration: undefined,
    mode: 0,
    runningRound: 0,
    runningRoundStart: undefined,
  };
}

export const GameState = {
  encode(message: GameState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(10).fork()).ldelim();
    }
    if (message.endTime !== undefined) {
      Timestamp.encode(toTimestamp(message.endTime), writer.uint32(18).fork()).ldelim();
    }
    if (message.totalRounds !== 0) {
      writer.uint32(24).uint32(message.totalRounds);
    }
    if (message.paused === true) {
      writer.uint32(32).bool(message.paused);
    }
    if (message.flagLifetimeRounds !== 0) {
      writer.uint32(40).uint32(message.flagLifetimeRounds);
    }
    if (message.roundDuration !== undefined) {
      Duration.encode(message.roundDuration, writer.uint32(50).fork()).ldelim();
    }
    if (message.mode !== 0) {
      writer.uint32(56).int32(message.mode);
    }
    if (message.runningRound !== 0) {
      writer.uint32(64).uint32(message.runningRound);
    }
    if (message.runningRoundStart !== undefined) {
      Timestamp.encode(toTimestamp(message.runningRoundStart), writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GameState {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGameState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.endTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.totalRounds = reader.uint32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.paused = reader.bool();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.flagLifetimeRounds = reader.uint32();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.roundDuration = Duration.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.mode = reader.int32() as any;
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.runningRound = reader.uint32();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.runningRoundStart = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<GameState, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<GameState | GameState[]> | Iterable<GameState | GameState[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [GameState.encode(p).finish()];
        }
      } else {
        yield* [GameState.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, GameState>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<GameState> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [GameState.decode(p)];
        }
      } else {
        yield* [GameState.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): GameState {
    return {
      startTime: isSet(object.startTime) ? fromJsonTimestamp(object.startTime) : undefined,
      endTime: isSet(object.endTime) ? fromJsonTimestamp(object.endTime) : undefined,
      totalRounds: isSet(object.totalRounds) ? Number(object.totalRounds) : 0,
      paused: isSet(object.paused) ? Boolean(object.paused) : false,
      flagLifetimeRounds: isSet(object.flagLifetimeRounds) ? Number(object.flagLifetimeRounds) : 0,
      roundDuration: isSet(object.roundDuration) ? Duration.fromJSON(object.roundDuration) : undefined,
      mode: isSet(object.mode) ? gameModeFromJSON(object.mode) : 0,
      runningRound: isSet(object.runningRound) ? Number(object.runningRound) : 0,
      runningRoundStart: isSet(object.runningRoundStart) ? fromJsonTimestamp(object.runningRoundStart) : undefined,
    };
  },

  toJSON(message: GameState): unknown {
    const obj: any = {};
    if (message.startTime !== undefined) {
      obj.startTime = message.startTime.toISOString();
    }
    if (message.endTime !== undefined) {
      obj.endTime = message.endTime.toISOString();
    }
    if (message.totalRounds !== 0) {
      obj.totalRounds = Math.round(message.totalRounds);
    }
    if (message.paused === true) {
      obj.paused = message.paused;
    }
    if (message.flagLifetimeRounds !== 0) {
      obj.flagLifetimeRounds = Math.round(message.flagLifetimeRounds);
    }
    if (message.roundDuration !== undefined) {
      obj.roundDuration = Duration.toJSON(message.roundDuration);
    }
    if (message.mode !== 0) {
      obj.mode = gameModeToJSON(message.mode);
    }
    if (message.runningRound !== 0) {
      obj.runningRound = Math.round(message.runningRound);
    }
    if (message.runningRoundStart !== undefined) {
      obj.runningRoundStart = message.runningRoundStart.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GameState>, I>>(base?: I): GameState {
    return GameState.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GameState>, I>>(object: I): GameState {
    const message = createBaseGameState();
    message.startTime = object.startTime ?? undefined;
    message.endTime = object.endTime ?? undefined;
    message.totalRounds = object.totalRounds ?? 0;
    message.paused = object.paused ?? false;
    message.flagLifetimeRounds = object.flagLifetimeRounds ?? 0;
    message.roundDuration = (object.roundDuration !== undefined && object.roundDuration !== null)
      ? Duration.fromPartial(object.roundDuration)
      : undefined;
    message.mode = object.mode ?? 0;
    message.runningRound = object.runningRound ?? 0;
    message.runningRoundStart = object.runningRoundStart ?? undefined;
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { version: undefined };
}

export const GetRequest = {
  encode(message: GetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.version = Version.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<GetRequest, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<GetRequest | GetRequest[]> | Iterable<GetRequest | GetRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [GetRequest.encode(p).finish()];
        }
      } else {
        yield* [GetRequest.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, GetRequest>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<GetRequest> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [GetRequest.decode(p)];
        }
      } else {
        yield* [GetRequest.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): GetRequest {
    return { version: isSet(object.version) ? Version.fromJSON(object.version) : undefined };
  },

  toJSON(message: GetRequest): unknown {
    const obj: any = {};
    if (message.version !== undefined) {
      obj.version = Version.toJSON(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRequest>, I>>(base?: I): GetRequest {
    return GetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(object: I): GetRequest {
    const message = createBaseGetRequest();
    message.version = (object.version !== undefined && object.version !== null)
      ? Version.fromPartial(object.version)
      : undefined;
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { gameState: undefined, version: undefined };
}

export const GetResponse = {
  encode(message: GetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.gameState !== undefined) {
      GameState.encode(message.gameState, writer.uint32(10).fork()).ldelim();
    }
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.gameState = GameState.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.version = Version.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<GetResponse, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<GetResponse | GetResponse[]> | Iterable<GetResponse | GetResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [GetResponse.encode(p).finish()];
        }
      } else {
        yield* [GetResponse.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, GetResponse>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<GetResponse> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [GetResponse.decode(p)];
        }
      } else {
        yield* [GetResponse.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): GetResponse {
    return {
      gameState: isSet(object.gameState) ? GameState.fromJSON(object.gameState) : undefined,
      version: isSet(object.version) ? Version.fromJSON(object.version) : undefined,
    };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    if (message.gameState !== undefined) {
      obj.gameState = GameState.toJSON(message.gameState);
    }
    if (message.version !== undefined) {
      obj.version = Version.toJSON(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetResponse>, I>>(base?: I): GetResponse {
    return GetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(object: I): GetResponse {
    const message = createBaseGetResponse();
    message.gameState = (object.gameState !== undefined && object.gameState !== null)
      ? GameState.fromPartial(object.gameState)
      : undefined;
    message.version = (object.version !== undefined && object.version !== null)
      ? Version.fromPartial(object.version)
      : undefined;
    return message;
  },
};

function createBaseUpdateRequest(): UpdateRequest {
  return {
    startTime: undefined,
    endTime: undefined,
    totalRounds: 0,
    paused: false,
    flagLifetimeRounds: 0,
    roundDuration: undefined,
    mode: 0,
  };
}

export const UpdateRequest = {
  encode(message: UpdateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(10).fork()).ldelim();
    }
    if (message.endTime !== undefined) {
      Timestamp.encode(toTimestamp(message.endTime), writer.uint32(18).fork()).ldelim();
    }
    if (message.totalRounds !== 0) {
      writer.uint32(24).uint32(message.totalRounds);
    }
    if (message.paused === true) {
      writer.uint32(32).bool(message.paused);
    }
    if (message.flagLifetimeRounds !== 0) {
      writer.uint32(40).uint32(message.flagLifetimeRounds);
    }
    if (message.roundDuration !== undefined) {
      Duration.encode(message.roundDuration, writer.uint32(50).fork()).ldelim();
    }
    if (message.mode !== 0) {
      writer.uint32(56).int32(message.mode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.endTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.totalRounds = reader.uint32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.paused = reader.bool();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.flagLifetimeRounds = reader.uint32();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.roundDuration = Duration.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.mode = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<UpdateRequest, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<UpdateRequest | UpdateRequest[]> | Iterable<UpdateRequest | UpdateRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateRequest.encode(p).finish()];
        }
      } else {
        yield* [UpdateRequest.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, UpdateRequest>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<UpdateRequest> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateRequest.decode(p)];
        }
      } else {
        yield* [UpdateRequest.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): UpdateRequest {
    return {
      startTime: isSet(object.startTime) ? fromJsonTimestamp(object.startTime) : undefined,
      endTime: isSet(object.endTime) ? fromJsonTimestamp(object.endTime) : undefined,
      totalRounds: isSet(object.totalRounds) ? Number(object.totalRounds) : 0,
      paused: isSet(object.paused) ? Boolean(object.paused) : false,
      flagLifetimeRounds: isSet(object.flagLifetimeRounds) ? Number(object.flagLifetimeRounds) : 0,
      roundDuration: isSet(object.roundDuration) ? Duration.fromJSON(object.roundDuration) : undefined,
      mode: isSet(object.mode) ? gameModeFromJSON(object.mode) : 0,
    };
  },

  toJSON(message: UpdateRequest): unknown {
    const obj: any = {};
    if (message.startTime !== undefined) {
      obj.startTime = message.startTime.toISOString();
    }
    if (message.endTime !== undefined) {
      obj.endTime = message.endTime.toISOString();
    }
    if (message.totalRounds !== 0) {
      obj.totalRounds = Math.round(message.totalRounds);
    }
    if (message.paused === true) {
      obj.paused = message.paused;
    }
    if (message.flagLifetimeRounds !== 0) {
      obj.flagLifetimeRounds = Math.round(message.flagLifetimeRounds);
    }
    if (message.roundDuration !== undefined) {
      obj.roundDuration = Duration.toJSON(message.roundDuration);
    }
    if (message.mode !== 0) {
      obj.mode = gameModeToJSON(message.mode);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRequest>, I>>(base?: I): UpdateRequest {
    return UpdateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRequest>, I>>(object: I): UpdateRequest {
    const message = createBaseUpdateRequest();
    message.startTime = object.startTime ?? undefined;
    message.endTime = object.endTime ?? undefined;
    message.totalRounds = object.totalRounds ?? 0;
    message.paused = object.paused ?? false;
    message.flagLifetimeRounds = object.flagLifetimeRounds ?? 0;
    message.roundDuration = (object.roundDuration !== undefined && object.roundDuration !== null)
      ? Duration.fromPartial(object.roundDuration)
      : undefined;
    message.mode = object.mode ?? 0;
    return message;
  },
};

function createBaseUpdateResponse(): UpdateResponse {
  return { gameState: undefined, version: undefined };
}

export const UpdateResponse = {
  encode(message: UpdateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.gameState !== undefined) {
      GameState.encode(message.gameState, writer.uint32(10).fork()).ldelim();
    }
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.gameState = GameState.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.version = Version.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<UpdateResponse, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<UpdateResponse | UpdateResponse[]> | Iterable<UpdateResponse | UpdateResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateResponse.encode(p).finish()];
        }
      } else {
        yield* [UpdateResponse.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, UpdateResponse>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<UpdateResponse> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateResponse.decode(p)];
        }
      } else {
        yield* [UpdateResponse.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): UpdateResponse {
    return {
      gameState: isSet(object.gameState) ? GameState.fromJSON(object.gameState) : undefined,
      version: isSet(object.version) ? Version.fromJSON(object.version) : undefined,
    };
  },

  toJSON(message: UpdateResponse): unknown {
    const obj: any = {};
    if (message.gameState !== undefined) {
      obj.gameState = GameState.toJSON(message.gameState);
    }
    if (message.version !== undefined) {
      obj.version = Version.toJSON(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateResponse>, I>>(base?: I): UpdateResponse {
    return UpdateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateResponse>, I>>(object: I): UpdateResponse {
    const message = createBaseUpdateResponse();
    message.gameState = (object.gameState !== undefined && object.gameState !== null)
      ? GameState.fromPartial(object.gameState)
      : undefined;
    message.version = (object.version !== undefined && object.version !== null)
      ? Version.fromPartial(object.version)
      : undefined;
    return message;
  },
};

function createBaseUpdateRoundRequest(): UpdateRoundRequest {
  return { runningRound: 0, runningRoundStart: undefined };
}

export const UpdateRoundRequest = {
  encode(message: UpdateRoundRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.runningRound !== 0) {
      writer.uint32(8).int32(message.runningRound);
    }
    if (message.runningRoundStart !== undefined) {
      Timestamp.encode(toTimestamp(message.runningRoundStart), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRoundRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRoundRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.runningRound = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.runningRoundStart = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<UpdateRoundRequest, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<UpdateRoundRequest | UpdateRoundRequest[]>
      | Iterable<UpdateRoundRequest | UpdateRoundRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateRoundRequest.encode(p).finish()];
        }
      } else {
        yield* [UpdateRoundRequest.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, UpdateRoundRequest>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<UpdateRoundRequest> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateRoundRequest.decode(p)];
        }
      } else {
        yield* [UpdateRoundRequest.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): UpdateRoundRequest {
    return {
      runningRound: isSet(object.runningRound) ? Number(object.runningRound) : 0,
      runningRoundStart: isSet(object.runningRoundStart) ? fromJsonTimestamp(object.runningRoundStart) : undefined,
    };
  },

  toJSON(message: UpdateRoundRequest): unknown {
    const obj: any = {};
    if (message.runningRound !== 0) {
      obj.runningRound = Math.round(message.runningRound);
    }
    if (message.runningRoundStart !== undefined) {
      obj.runningRoundStart = message.runningRoundStart.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRoundRequest>, I>>(base?: I): UpdateRoundRequest {
    return UpdateRoundRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRoundRequest>, I>>(object: I): UpdateRoundRequest {
    const message = createBaseUpdateRoundRequest();
    message.runningRound = object.runningRound ?? 0;
    message.runningRoundStart = object.runningRoundStart ?? undefined;
    return message;
  },
};

function createBaseUpdateRoundResponse(): UpdateRoundResponse {
  return { gameState: undefined, version: undefined };
}

export const UpdateRoundResponse = {
  encode(message: UpdateRoundResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.gameState !== undefined) {
      GameState.encode(message.gameState, writer.uint32(10).fork()).ldelim();
    }
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRoundResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRoundResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.gameState = GameState.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.version = Version.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<UpdateRoundResponse, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<UpdateRoundResponse | UpdateRoundResponse[]>
      | Iterable<UpdateRoundResponse | UpdateRoundResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateRoundResponse.encode(p).finish()];
        }
      } else {
        yield* [UpdateRoundResponse.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, UpdateRoundResponse>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<UpdateRoundResponse> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [UpdateRoundResponse.decode(p)];
        }
      } else {
        yield* [UpdateRoundResponse.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): UpdateRoundResponse {
    return {
      gameState: isSet(object.gameState) ? GameState.fromJSON(object.gameState) : undefined,
      version: isSet(object.version) ? Version.fromJSON(object.version) : undefined,
    };
  },

  toJSON(message: UpdateRoundResponse): unknown {
    const obj: any = {};
    if (message.gameState !== undefined) {
      obj.gameState = GameState.toJSON(message.gameState);
    }
    if (message.version !== undefined) {
      obj.version = Version.toJSON(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRoundResponse>, I>>(base?: I): UpdateRoundResponse {
    return UpdateRoundResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRoundResponse>, I>>(object: I): UpdateRoundResponse {
    const message = createBaseUpdateRoundResponse();
    message.gameState = (object.gameState !== undefined && object.gameState !== null)
      ? GameState.fromPartial(object.gameState)
      : undefined;
    message.version = (object.version !== undefined && object.version !== null)
      ? Version.fromPartial(object.version)
      : undefined;
    return message;
  },
};

declare const self: any | undefined;
declare const window: any | undefined;
declare const global: any | undefined;
const tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | bigint | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function toTimestamp(date: Date): Timestamp {
  const seconds = BigInt(Math.trunc(date.getTime() / 1_000));
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = (Number(t.seconds.toString()) || 0) * 1_000;
  millis += (t.nanos || 0) / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
