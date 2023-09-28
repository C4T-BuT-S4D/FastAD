/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "data.teams";

export interface Team {
  id: number;
  name: string;
  address: string;
  token: string;
  labels: { [key: string]: string };
}

export interface Team_LabelsEntry {
  key: string;
  value: string;
}

export interface ListRequest {
  lastUpdate: bigint;
}

export interface ListResponse {
  teams: Team[];
  lastUpdate: bigint;
}

export interface CreateBatchRequest {
  teams: Team[];
}

export interface CreateBatchResponse {
  teams: Team[];
}

function createBaseTeam(): Team {
  return {id: 0, name: "", address: "", token: "", labels: {}};
}

export const Team = {
  encode(message: Team, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).int32(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.address !== "") {
      writer.uint32(26).string(message.address);
    }
    if (message.token !== "") {
      writer.uint32(34).string(message.token);
    }
    Object.entries(message.labels).forEach(([key, value]) => {
      Team_LabelsEntry.encode({key: key as any, value}, writer.uint32(42).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Team {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTeam();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.address = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.token = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          const entry5 = Team_LabelsEntry.decode(reader, reader.uint32());
          if (entry5.value !== undefined) {
            message.labels[entry5.key] = entry5.value;
          }
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
  // Transform<Team, Uint8Array>
  async* encodeTransform(source: AsyncIterable<Team | Team[]> | Iterable<Team | Team[]>): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Team.encode(p).finish()];
        }
      } else {
        yield* [Team.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Team>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Team> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Team.decode(p)];
        }
      } else {
        yield* [Team.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): Team {
    return {
      id: isSet(object.id) ? Number(object.id) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      address: isSet(object.address) ? String(object.address) : "",
      token: isSet(object.token) ? String(object.token) : "",
      labels: isObject(object.labels)
          ? Object.entries(object.labels).reduce<{ [key: string]: string }>((acc, [key, value]) => {
            acc[key] = String(value);
            return acc;
          }, {})
          : {},
    };
  },

  toJSON(message: Team): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.address !== "") {
      obj.address = message.address;
    }
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.labels) {
      const entries = Object.entries(message.labels);
      if (entries.length > 0) {
        obj.labels = {};
        entries.forEach(([k, v]) => {
          obj.labels[k] = v;
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Team>, I>>(base?: I): Team {
    return Team.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Team>, I>>(object: I): Team {
    const message = createBaseTeam();
    message.id = object.id ?? 0;
    message.name = object.name ?? "";
    message.address = object.address ?? "";
    message.token = object.token ?? "";
    message.labels = Object.entries(object.labels ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = String(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseTeam_LabelsEntry(): Team_LabelsEntry {
  return {key: "", value: ""};
}

export const Team_LabelsEntry = {
  encode(message: Team_LabelsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Team_LabelsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTeam_LabelsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = reader.string();
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
  // Transform<Team_LabelsEntry, Uint8Array>
  async* encodeTransform(
      source: AsyncIterable<Team_LabelsEntry | Team_LabelsEntry[]> | Iterable<Team_LabelsEntry | Team_LabelsEntry[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Team_LabelsEntry.encode(p).finish()];
        }
      } else {
        yield* [Team_LabelsEntry.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Team_LabelsEntry>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Team_LabelsEntry> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Team_LabelsEntry.decode(p)];
        }
      } else {
        yield* [Team_LabelsEntry.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): Team_LabelsEntry {
    return {key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : ""};
  },

  toJSON(message: Team_LabelsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== "") {
      obj.value = message.value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Team_LabelsEntry>, I>>(base?: I): Team_LabelsEntry {
    return Team_LabelsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Team_LabelsEntry>, I>>(object: I): Team_LabelsEntry {
    const message = createBaseTeam_LabelsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseListRequest(): ListRequest {
  return {lastUpdate: BigInt("0")};
}

export const ListRequest = {
  encode(message: ListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.lastUpdate !== BigInt("0")) {
      writer.uint32(8).int64(message.lastUpdate.toString());
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.lastUpdate = longToBigint(reader.int64() as Long);
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
  // Transform<ListRequest, Uint8Array>
  async* encodeTransform(
      source: AsyncIterable<ListRequest | ListRequest[]> | Iterable<ListRequest | ListRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [ListRequest.encode(p).finish()];
        }
      } else {
        yield* [ListRequest.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ListRequest>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ListRequest> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [ListRequest.decode(p)];
        }
      } else {
        yield* [ListRequest.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): ListRequest {
    return {lastUpdate: isSet(object.lastUpdate) ? BigInt(object.lastUpdate) : BigInt("0")};
  },

  toJSON(message: ListRequest): unknown {
    const obj: any = {};
    if (message.lastUpdate !== BigInt("0")) {
      obj.lastUpdate = message.lastUpdate.toString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListRequest>, I>>(base?: I): ListRequest {
    return ListRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRequest>, I>>(object: I): ListRequest {
    const message = createBaseListRequest();
    message.lastUpdate = object.lastUpdate ?? BigInt("0");
    return message;
  },
};

function createBaseListResponse(): ListResponse {
  return {teams: [], lastUpdate: BigInt("0")};
}

export const ListResponse = {
  encode(message: ListResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.teams) {
      Team.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.lastUpdate !== BigInt("0")) {
      writer.uint32(16).int64(message.lastUpdate.toString());
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.teams.push(Team.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.lastUpdate = longToBigint(reader.int64() as Long);
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
  // Transform<ListResponse, Uint8Array>
  async* encodeTransform(
      source: AsyncIterable<ListResponse | ListResponse[]> | Iterable<ListResponse | ListResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [ListResponse.encode(p).finish()];
        }
      } else {
        yield* [ListResponse.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ListResponse>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ListResponse> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [ListResponse.decode(p)];
        }
      } else {
        yield* [ListResponse.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): ListResponse {
    return {
      teams: Array.isArray(object?.teams) ? object.teams.map((e: any) => Team.fromJSON(e)) : [],
      lastUpdate: isSet(object.lastUpdate) ? BigInt(object.lastUpdate) : BigInt("0"),
    };
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    if (message.teams?.length) {
      obj.teams = message.teams.map((e) => Team.toJSON(e));
    }
    if (message.lastUpdate !== BigInt("0")) {
      obj.lastUpdate = message.lastUpdate.toString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListResponse>, I>>(base?: I): ListResponse {
    return ListResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListResponse>, I>>(object: I): ListResponse {
    const message = createBaseListResponse();
    message.teams = object.teams?.map((e) => Team.fromPartial(e)) || [];
    message.lastUpdate = object.lastUpdate ?? BigInt("0");
    return message;
  },
};

function createBaseCreateBatchRequest(): CreateBatchRequest {
  return {teams: []};
}

export const CreateBatchRequest = {
  encode(message: CreateBatchRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.teams) {
      Team.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateBatchRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateBatchRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.teams.push(Team.decode(reader, reader.uint32()));
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
  // Transform<CreateBatchRequest, Uint8Array>
  async* encodeTransform(
      source:
          | AsyncIterable<CreateBatchRequest | CreateBatchRequest[]>
          | Iterable<CreateBatchRequest | CreateBatchRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [CreateBatchRequest.encode(p).finish()];
        }
      } else {
        yield* [CreateBatchRequest.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, CreateBatchRequest>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<CreateBatchRequest> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [CreateBatchRequest.decode(p)];
        }
      } else {
        yield* [CreateBatchRequest.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): CreateBatchRequest {
    return {teams: Array.isArray(object?.teams) ? object.teams.map((e: any) => Team.fromJSON(e)) : []};
  },

  toJSON(message: CreateBatchRequest): unknown {
    const obj: any = {};
    if (message.teams?.length) {
      obj.teams = message.teams.map((e) => Team.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateBatchRequest>, I>>(base?: I): CreateBatchRequest {
    return CreateBatchRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateBatchRequest>, I>>(object: I): CreateBatchRequest {
    const message = createBaseCreateBatchRequest();
    message.teams = object.teams?.map((e) => Team.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCreateBatchResponse(): CreateBatchResponse {
  return {teams: []};
}

export const CreateBatchResponse = {
  encode(message: CreateBatchResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.teams) {
      Team.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateBatchResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateBatchResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.teams.push(Team.decode(reader, reader.uint32()));
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
  // Transform<CreateBatchResponse, Uint8Array>
  async* encodeTransform(
      source:
          | AsyncIterable<CreateBatchResponse | CreateBatchResponse[]>
          | Iterable<CreateBatchResponse | CreateBatchResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [CreateBatchResponse.encode(p).finish()];
        }
      } else {
        yield* [CreateBatchResponse.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, CreateBatchResponse>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<CreateBatchResponse> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [CreateBatchResponse.decode(p)];
        }
      } else {
        yield* [CreateBatchResponse.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): CreateBatchResponse {
    return {teams: Array.isArray(object?.teams) ? object.teams.map((e: any) => Team.fromJSON(e)) : []};
  },

  toJSON(message: CreateBatchResponse): unknown {
    const obj: any = {};
    if (message.teams?.length) {
      obj.teams = message.teams.map((e) => Team.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateBatchResponse>, I>>(base?: I): CreateBatchResponse {
    return CreateBatchResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateBatchResponse>, I>>(object: I): CreateBatchResponse {
    const message = createBaseCreateBatchResponse();
    message.teams = object.teams?.map((e) => Team.fromPartial(e)) || [];
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | bigint | undefined;

export type DeepPartial<T> = T extends Builtin ? T
    : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
        : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
            : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
    : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToBigint(long: Long) {
  return BigInt(long.toString());
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
