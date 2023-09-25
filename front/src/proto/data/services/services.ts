/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import {Type, typeFromJSON, typeToJSON} from "../../checker/checker";

export const protobufPackage = "data.services";

export interface Service {
  id: number;
  name: string;
  checker: Service_Checker | undefined;
  defaultScore: number;
  gets: number;
  puts: number;
}

export interface Service_Checker {
  type: Type;
  path: string;
  timeoutSeconds: number;
}

export interface ListRequest {
  lastUpdate: number;
}

export interface ListResponse {
  services: Service[];
  lastUpdate: number;
}

function createBaseService(): Service {
  return {id: 0, name: "", checker: undefined, defaultScore: 0, gets: 0, puts: 0};
}

export const Service = {
  encode(message: Service, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).int32(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.checker !== undefined) {
      Service_Checker.encode(message.checker, writer.uint32(26).fork()).ldelim();
    }
    if (message.defaultScore !== 0) {
      writer.uint32(33).double(message.defaultScore);
    }
    if (message.gets !== 0) {
      writer.uint32(40).int32(message.gets);
    }
    if (message.puts !== 0) {
      writer.uint32(48).int32(message.puts);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Service {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseService();
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

          message.checker = Service_Checker.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.defaultScore = reader.double();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.gets = reader.int32();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.puts = reader.int32();
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
  // Transform<Service, Uint8Array>
  async* encodeTransform(
      source: AsyncIterable<Service | Service[]> | Iterable<Service | Service[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service.encode(p).finish()];
        }
      } else {
        yield* [Service.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Service>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Service> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service.decode(p)];
        }
      } else {
        yield* [Service.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): Service {
    return {
      id: isSet(object.id) ? Number(object.id) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      checker: isSet(object.checker) ? Service_Checker.fromJSON(object.checker) : undefined,
      defaultScore: isSet(object.defaultScore) ? Number(object.defaultScore) : 0,
      gets: isSet(object.gets) ? Number(object.gets) : 0,
      puts: isSet(object.puts) ? Number(object.puts) : 0,
    };
  },

  toJSON(message: Service): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.checker !== undefined) {
      obj.checker = Service_Checker.toJSON(message.checker);
    }
    if (message.defaultScore !== 0) {
      obj.defaultScore = message.defaultScore;
    }
    if (message.gets !== 0) {
      obj.gets = Math.round(message.gets);
    }
    if (message.puts !== 0) {
      obj.puts = Math.round(message.puts);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Service>, I>>(base?: I): Service {
    return Service.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Service>, I>>(object: I): Service {
    const message = createBaseService();
    message.id = object.id ?? 0;
    message.name = object.name ?? "";
    message.checker = (object.checker !== undefined && object.checker !== null)
        ? Service_Checker.fromPartial(object.checker)
        : undefined;
    message.defaultScore = object.defaultScore ?? 0;
    message.gets = object.gets ?? 0;
    message.puts = object.puts ?? 0;
    return message;
  },
};

function createBaseService_Checker(): Service_Checker {
  return {type: 0, path: "", timeoutSeconds: 0};
}

export const Service_Checker = {
  encode(message: Service_Checker, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    if (message.path !== "") {
      writer.uint32(18).string(message.path);
    }
    if (message.timeoutSeconds !== 0) {
      writer.uint32(24).int32(message.timeoutSeconds);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Service_Checker {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseService_Checker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.path = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.timeoutSeconds = reader.int32();
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
  // Transform<Service_Checker, Uint8Array>
  async* encodeTransform(
      source: AsyncIterable<Service_Checker | Service_Checker[]> | Iterable<Service_Checker | Service_Checker[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service_Checker.encode(p).finish()];
        }
      } else {
        yield* [Service_Checker.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Service_Checker>
  async* decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Service_Checker> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service_Checker.decode(p)];
        }
      } else {
        yield* [Service_Checker.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): Service_Checker {
    return {
      type: isSet(object.type) ? typeFromJSON(object.type) : 0,
      path: isSet(object.path) ? String(object.path) : "",
      timeoutSeconds: isSet(object.timeoutSeconds) ? Number(object.timeoutSeconds) : 0,
    };
  },

  toJSON(message: Service_Checker): unknown {
    const obj: any = {};
    if (message.type !== 0) {
      obj.type = typeToJSON(message.type);
    }
    if (message.path !== "") {
      obj.path = message.path;
    }
    if (message.timeoutSeconds !== 0) {
      obj.timeoutSeconds = Math.round(message.timeoutSeconds);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Service_Checker>, I>>(base?: I): Service_Checker {
    return Service_Checker.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Service_Checker>, I>>(object: I): Service_Checker {
    const message = createBaseService_Checker();
    message.type = object.type ?? 0;
    message.path = object.path ?? "";
    message.timeoutSeconds = object.timeoutSeconds ?? 0;
    return message;
  },
};

function createBaseListRequest(): ListRequest {
  return {lastUpdate: 0};
}

export const ListRequest = {
  encode(message: ListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.lastUpdate !== 0) {
      writer.uint32(8).int64(message.lastUpdate);
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

          message.lastUpdate = longToNumber(reader.int64() as Long);
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
    return {lastUpdate: isSet(object.lastUpdate) ? Number(object.lastUpdate) : 0};
  },

  toJSON(message: ListRequest): unknown {
    const obj: any = {};
    if (message.lastUpdate !== 0) {
      obj.lastUpdate = Math.round(message.lastUpdate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListRequest>, I>>(base?: I): ListRequest {
    return ListRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRequest>, I>>(object: I): ListRequest {
    const message = createBaseListRequest();
    message.lastUpdate = object.lastUpdate ?? 0;
    return message;
  },
};

function createBaseListResponse(): ListResponse {
  return {services: [], lastUpdate: 0};
}

export const ListResponse = {
  encode(message: ListResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.services) {
      Service.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.lastUpdate !== 0) {
      writer.uint32(16).int64(message.lastUpdate);
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

          message.services.push(Service.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.lastUpdate = longToNumber(reader.int64() as Long);
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
      services: Array.isArray(object?.services) ? object.services.map((e: any) => Service.fromJSON(e)) : [],
      lastUpdate: isSet(object.lastUpdate) ? Number(object.lastUpdate) : 0,
    };
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    if (message.services?.length) {
      obj.services = message.services.map((e) => Service.toJSON(e));
    }
    if (message.lastUpdate !== 0) {
      obj.lastUpdate = Math.round(message.lastUpdate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListResponse>, I>>(base?: I): ListResponse {
    return ListResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListResponse>, I>>(object: I): ListResponse {
    const message = createBaseListResponse();
    message.services = object.services?.map((e) => Service.fromPartial(e)) || [];
    message.lastUpdate = object.lastUpdate ?? 0;
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
    : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
        : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
            : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
    : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
