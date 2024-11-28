// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.3.0
//   protoc               unknown
// source: data/teams/teams.proto

/* eslint-disable */
import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import { Version } from "../version/version";

export const protobufPackage = "data.teams";

export interface Team {
  id: bigint;
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
  version: Version | undefined;
}

export interface ListResponse {
  teams: Team[];
  version: Version | undefined;
}

export interface CreateBatchRequest {
  teams: Team[];
}

export interface CreateBatchResponse {
  teams: Team[];
}

function createBaseTeam(): Team {
  return { id: 0n, name: "", address: "", token: "", labels: {} };
}

export const Team: MessageFns<Team> = {
  encode(message: Team, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.id !== 0n) {
      if (BigInt.asIntN(64, message.id) !== message.id) {
        throw new globalThis.Error("value provided for field message.id of type int64 too large");
      }
      writer.uint32(8).int64(message.id);
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
      Team_LabelsEntry.encode({ key: key as any, value }, writer.uint32(42).fork()).join();
    });
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): Team {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTeam();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 8) {
            break;
          }

          message.id = reader.int64() as bigint;
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
        }
        case 3: {
          if (tag !== 26) {
            break;
          }

          message.address = reader.string();
          continue;
        }
        case 4: {
          if (tag !== 34) {
            break;
          }

          message.token = reader.string();
          continue;
        }
        case 5: {
          if (tag !== 42) {
            break;
          }

          const entry5 = Team_LabelsEntry.decode(reader, reader.uint32());
          if (entry5.value !== undefined) {
            message.labels[entry5.key] = entry5.value;
          }
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<Team, Uint8Array>
  async *encodeTransform(source: AsyncIterable<Team | Team[]> | Iterable<Team | Team[]>): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [Team.encode(p).finish()];
        }
      } else {
        yield* [Team.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Team>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Team> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [Team.decode(p)];
        }
      } else {
        yield* [Team.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): Team {
    return {
      id: isSet(object.id) ? BigInt(object.id) : 0n,
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      address: isSet(object.address) ? globalThis.String(object.address) : "",
      token: isSet(object.token) ? globalThis.String(object.token) : "",
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
    if (message.id !== 0n) {
      obj.id = message.id.toString();
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
    message.id = object.id ?? 0n;
    message.name = object.name ?? "";
    message.address = object.address ?? "";
    message.token = object.token ?? "";
    message.labels = Object.entries(object.labels ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = globalThis.String(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseTeam_LabelsEntry(): Team_LabelsEntry {
  return { key: "", value: "" };
}

export const Team_LabelsEntry: MessageFns<Team_LabelsEntry> = {
  encode(message: Team_LabelsEntry, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): Team_LabelsEntry {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTeam_LabelsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.value = reader.string();
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<Team_LabelsEntry, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<Team_LabelsEntry | Team_LabelsEntry[]> | Iterable<Team_LabelsEntry | Team_LabelsEntry[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [Team_LabelsEntry.encode(p).finish()];
        }
      } else {
        yield* [Team_LabelsEntry.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Team_LabelsEntry>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Team_LabelsEntry> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [Team_LabelsEntry.decode(p)];
        }
      } else {
        yield* [Team_LabelsEntry.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): Team_LabelsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? globalThis.String(object.value) : "",
    };
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
  return { version: undefined };
}

export const ListRequest: MessageFns<ListRequest> = {
  encode(message: ListRequest, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): ListRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.version = Version.decode(reader, reader.uint32());
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<ListRequest, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<ListRequest | ListRequest[]> | Iterable<ListRequest | ListRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ListRequest.encode(p).finish()];
        }
      } else {
        yield* [ListRequest.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ListRequest>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ListRequest> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ListRequest.decode(p)];
        }
      } else {
        yield* [ListRequest.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): ListRequest {
    return { version: isSet(object.version) ? Version.fromJSON(object.version) : undefined };
  },

  toJSON(message: ListRequest): unknown {
    const obj: any = {};
    if (message.version !== undefined) {
      obj.version = Version.toJSON(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListRequest>, I>>(base?: I): ListRequest {
    return ListRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRequest>, I>>(object: I): ListRequest {
    const message = createBaseListRequest();
    message.version = (object.version !== undefined && object.version !== null)
      ? Version.fromPartial(object.version)
      : undefined;
    return message;
  },
};

function createBaseListResponse(): ListResponse {
  return { teams: [], version: undefined };
}

export const ListResponse: MessageFns<ListResponse> = {
  encode(message: ListResponse, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    for (const v of message.teams) {
      Team.encode(v!, writer.uint32(10).fork()).join();
    }
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(18).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): ListResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.teams.push(Team.decode(reader, reader.uint32()));
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.version = Version.decode(reader, reader.uint32());
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<ListResponse, Uint8Array>
  async *encodeTransform(
    source: AsyncIterable<ListResponse | ListResponse[]> | Iterable<ListResponse | ListResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ListResponse.encode(p).finish()];
        }
      } else {
        yield* [ListResponse.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ListResponse>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ListResponse> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ListResponse.decode(p)];
        }
      } else {
        yield* [ListResponse.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): ListResponse {
    return {
      teams: globalThis.Array.isArray(object?.teams) ? object.teams.map((e: any) => Team.fromJSON(e)) : [],
      version: isSet(object.version) ? Version.fromJSON(object.version) : undefined,
    };
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    if (message.teams?.length) {
      obj.teams = message.teams.map((e) => Team.toJSON(e));
    }
    if (message.version !== undefined) {
      obj.version = Version.toJSON(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListResponse>, I>>(base?: I): ListResponse {
    return ListResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListResponse>, I>>(object: I): ListResponse {
    const message = createBaseListResponse();
    message.teams = object.teams?.map((e) => Team.fromPartial(e)) || [];
    message.version = (object.version !== undefined && object.version !== null)
      ? Version.fromPartial(object.version)
      : undefined;
    return message;
  },
};

function createBaseCreateBatchRequest(): CreateBatchRequest {
  return { teams: [] };
}

export const CreateBatchRequest: MessageFns<CreateBatchRequest> = {
  encode(message: CreateBatchRequest, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    for (const v of message.teams) {
      Team.encode(v!, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): CreateBatchRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateBatchRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.teams.push(Team.decode(reader, reader.uint32()));
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<CreateBatchRequest, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<CreateBatchRequest | CreateBatchRequest[]>
      | Iterable<CreateBatchRequest | CreateBatchRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [CreateBatchRequest.encode(p).finish()];
        }
      } else {
        yield* [CreateBatchRequest.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, CreateBatchRequest>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<CreateBatchRequest> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [CreateBatchRequest.decode(p)];
        }
      } else {
        yield* [CreateBatchRequest.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): CreateBatchRequest {
    return { teams: globalThis.Array.isArray(object?.teams) ? object.teams.map((e: any) => Team.fromJSON(e)) : [] };
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
  return { teams: [] };
}

export const CreateBatchResponse: MessageFns<CreateBatchResponse> = {
  encode(message: CreateBatchResponse, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    for (const v of message.teams) {
      Team.encode(v!, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): CreateBatchResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateBatchResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.teams.push(Team.decode(reader, reader.uint32()));
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<CreateBatchResponse, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<CreateBatchResponse | CreateBatchResponse[]>
      | Iterable<CreateBatchResponse | CreateBatchResponse[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [CreateBatchResponse.encode(p).finish()];
        }
      } else {
        yield* [CreateBatchResponse.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, CreateBatchResponse>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<CreateBatchResponse> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [CreateBatchResponse.decode(p)];
        }
      } else {
        yield* [CreateBatchResponse.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): CreateBatchResponse {
    return { teams: globalThis.Array.isArray(object?.teams) ? object.teams.map((e: any) => Team.fromJSON(e)) : [] };
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
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export interface MessageFns<T> {
  encode(message: T, writer?: BinaryWriter): BinaryWriter;
  decode(input: BinaryReader | Uint8Array, length?: number): T;
  encodeTransform(source: AsyncIterable<T | T[]> | Iterable<T | T[]>): AsyncIterable<Uint8Array>;
  decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<T>;
  fromJSON(object: any): T;
  toJSON(message: T): unknown;
  create<I extends Exact<DeepPartial<T>, I>>(base?: I): T;
  fromPartial<I extends Exact<DeepPartial<T>, I>>(object: I): T;
}
