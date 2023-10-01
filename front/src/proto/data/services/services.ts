/* eslint-disable */
import _m0 from "protobufjs/minimal";
import {Action, actionFromJSON, actionToJSON, Type, typeFromJSON, typeToJSON} from "../../checker/checker";
import {Duration} from "../../google/protobuf/duration";
import {Version} from "../version/version";

export const protobufPackage = "data.services";

export interface Service {
  id: number;
  name: string;
  checker: Service_Checker | undefined;
  defaultScore: number;
}

export interface Service_Checker {
  type: Type;
  path: string;
  defaultTimeout: Duration | undefined;
  actionTimeouts: Service_Checker_ActionTimeout[];
  actionRunCounts: Service_Checker_ActionRunCount[];
}

export interface Service_Checker_ActionTimeout {
  action: Action;
  timeout: Duration | undefined;
}

export interface Service_Checker_ActionRunCount {
  action: Action;
  runCount: number;
}

export interface ListRequest {
  version: Version | undefined;
}

export interface ListResponse {
  services: Service[];
  version: Version | undefined;
}

export interface CreateBatchRequest {
  services: Service[];
}

export interface CreateBatchResponse {
  services: Service[];
}

function createBaseService(): Service {
  return { id: 0, name: "", checker: undefined, defaultScore: 0 };
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
  async *encodeTransform(
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
  async *decodeTransform(
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
    return message;
  },
};

function createBaseService_Checker(): Service_Checker {
  return { type: 0, path: "", defaultTimeout: undefined, actionTimeouts: [], actionRunCounts: [] };
}

export const Service_Checker = {
  encode(message: Service_Checker, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    if (message.path !== "") {
      writer.uint32(18).string(message.path);
    }
    if (message.defaultTimeout !== undefined) {
      Duration.encode(message.defaultTimeout, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.actionTimeouts) {
      Service_Checker_ActionTimeout.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.actionRunCounts) {
      Service_Checker_ActionRunCount.encode(v!, writer.uint32(42).fork()).ldelim();
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
          if (tag !== 26) {
            break;
          }

          message.defaultTimeout = Duration.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.actionTimeouts.push(Service_Checker_ActionTimeout.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.actionRunCounts.push(Service_Checker_ActionRunCount.decode(reader, reader.uint32()));
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
  async *encodeTransform(
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
  async *decodeTransform(
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
      defaultTimeout: isSet(object.defaultTimeout) ? Duration.fromJSON(object.defaultTimeout) : undefined,
      actionTimeouts: Array.isArray(object?.actionTimeouts)
        ? object.actionTimeouts.map((e: any) => Service_Checker_ActionTimeout.fromJSON(e))
        : [],
      actionRunCounts: Array.isArray(object?.actionRunCounts)
        ? object.actionRunCounts.map((e: any) => Service_Checker_ActionRunCount.fromJSON(e))
        : [],
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
    if (message.defaultTimeout !== undefined) {
      obj.defaultTimeout = Duration.toJSON(message.defaultTimeout);
    }
    if (message.actionTimeouts?.length) {
      obj.actionTimeouts = message.actionTimeouts.map((e) => Service_Checker_ActionTimeout.toJSON(e));
    }
    if (message.actionRunCounts?.length) {
      obj.actionRunCounts = message.actionRunCounts.map((e) => Service_Checker_ActionRunCount.toJSON(e));
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
    message.defaultTimeout = (object.defaultTimeout !== undefined && object.defaultTimeout !== null)
      ? Duration.fromPartial(object.defaultTimeout)
      : undefined;
    message.actionTimeouts = object.actionTimeouts?.map((e) => Service_Checker_ActionTimeout.fromPartial(e)) || [];
    message.actionRunCounts = object.actionRunCounts?.map((e) => Service_Checker_ActionRunCount.fromPartial(e)) || [];
    return message;
  },
};

function createBaseService_Checker_ActionTimeout(): Service_Checker_ActionTimeout {
  return { action: 0, timeout: undefined };
}

export const Service_Checker_ActionTimeout = {
  encode(message: Service_Checker_ActionTimeout, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.action !== 0) {
      writer.uint32(8).int32(message.action);
    }
    if (message.timeout !== undefined) {
      Duration.encode(message.timeout, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Service_Checker_ActionTimeout {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseService_Checker_ActionTimeout();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.action = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.timeout = Duration.decode(reader, reader.uint32());
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
  // Transform<Service_Checker_ActionTimeout, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<Service_Checker_ActionTimeout | Service_Checker_ActionTimeout[]>
      | Iterable<Service_Checker_ActionTimeout | Service_Checker_ActionTimeout[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service_Checker_ActionTimeout.encode(p).finish()];
        }
      } else {
        yield* [Service_Checker_ActionTimeout.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Service_Checker_ActionTimeout>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Service_Checker_ActionTimeout> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service_Checker_ActionTimeout.decode(p)];
        }
      } else {
        yield* [Service_Checker_ActionTimeout.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): Service_Checker_ActionTimeout {
    return {
      action: isSet(object.action) ? actionFromJSON(object.action) : 0,
      timeout: isSet(object.timeout) ? Duration.fromJSON(object.timeout) : undefined,
    };
  },

  toJSON(message: Service_Checker_ActionTimeout): unknown {
    const obj: any = {};
    if (message.action !== 0) {
      obj.action = actionToJSON(message.action);
    }
    if (message.timeout !== undefined) {
      obj.timeout = Duration.toJSON(message.timeout);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Service_Checker_ActionTimeout>, I>>(base?: I): Service_Checker_ActionTimeout {
    return Service_Checker_ActionTimeout.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Service_Checker_ActionTimeout>, I>>(
    object: I,
  ): Service_Checker_ActionTimeout {
    const message = createBaseService_Checker_ActionTimeout();
    message.action = object.action ?? 0;
    message.timeout = (object.timeout !== undefined && object.timeout !== null)
      ? Duration.fromPartial(object.timeout)
      : undefined;
    return message;
  },
};

function createBaseService_Checker_ActionRunCount(): Service_Checker_ActionRunCount {
  return { action: 0, runCount: 0 };
}

export const Service_Checker_ActionRunCount = {
  encode(message: Service_Checker_ActionRunCount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.action !== 0) {
      writer.uint32(8).int32(message.action);
    }
    if (message.runCount !== 0) {
      writer.uint32(16).int32(message.runCount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Service_Checker_ActionRunCount {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseService_Checker_ActionRunCount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.action = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.runCount = reader.int32();
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
  // Transform<Service_Checker_ActionRunCount, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<Service_Checker_ActionRunCount | Service_Checker_ActionRunCount[]>
      | Iterable<Service_Checker_ActionRunCount | Service_Checker_ActionRunCount[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service_Checker_ActionRunCount.encode(p).finish()];
        }
      } else {
        yield* [Service_Checker_ActionRunCount.encode(pkt).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, Service_Checker_ActionRunCount>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<Service_Checker_ActionRunCount> {
    for await (const pkt of source) {
      if (Array.isArray(pkt)) {
        for (const p of pkt) {
          yield* [Service_Checker_ActionRunCount.decode(p)];
        }
      } else {
        yield* [Service_Checker_ActionRunCount.decode(pkt)];
      }
    }
  },

  fromJSON(object: any): Service_Checker_ActionRunCount {
    return {
      action: isSet(object.action) ? actionFromJSON(object.action) : 0,
      runCount: isSet(object.runCount) ? Number(object.runCount) : 0,
    };
  },

  toJSON(message: Service_Checker_ActionRunCount): unknown {
    const obj: any = {};
    if (message.action !== 0) {
      obj.action = actionToJSON(message.action);
    }
    if (message.runCount !== 0) {
      obj.runCount = Math.round(message.runCount);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Service_Checker_ActionRunCount>, I>>(base?: I): Service_Checker_ActionRunCount {
    return Service_Checker_ActionRunCount.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Service_Checker_ActionRunCount>, I>>(
    object: I,
  ): Service_Checker_ActionRunCount {
    const message = createBaseService_Checker_ActionRunCount();
    message.action = object.action ?? 0;
    message.runCount = object.runCount ?? 0;
    return message;
  },
};

function createBaseListRequest(): ListRequest {
  return {version: undefined};
}

export const ListRequest = {
  encode(message: ListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(10).fork()).ldelim();
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
  // Transform<ListRequest, Uint8Array>
  async *encodeTransform(
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
  async *decodeTransform(
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
    return {version: isSet(object.version) ? Version.fromJSON(object.version) : undefined};
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
  return {services: [], version: undefined};
}

export const ListResponse = {
  encode(message: ListResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.services) {
      Service.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.version !== undefined) {
      Version.encode(message.version, writer.uint32(18).fork()).ldelim();
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
  // Transform<ListResponse, Uint8Array>
  async *encodeTransform(
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
  async *decodeTransform(
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
      version: isSet(object.version) ? Version.fromJSON(object.version) : undefined,
    };
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    if (message.services?.length) {
      obj.services = message.services.map((e) => Service.toJSON(e));
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
    message.services = object.services?.map((e) => Service.fromPartial(e)) || [];
    message.version = (object.version !== undefined && object.version !== null)
        ? Version.fromPartial(object.version)
        : undefined;
    return message;
  },
};

function createBaseCreateBatchRequest(): CreateBatchRequest {
  return { services: [] };
}

export const CreateBatchRequest = {
  encode(message: CreateBatchRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.services) {
      Service.encode(v!, writer.uint32(10).fork()).ldelim();
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

          message.services.push(Service.decode(reader, reader.uint32()));
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
  async *encodeTransform(
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
  async *decodeTransform(
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
    return { services: Array.isArray(object?.services) ? object.services.map((e: any) => Service.fromJSON(e)) : [] };
  },

  toJSON(message: CreateBatchRequest): unknown {
    const obj: any = {};
    if (message.services?.length) {
      obj.services = message.services.map((e) => Service.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateBatchRequest>, I>>(base?: I): CreateBatchRequest {
    return CreateBatchRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateBatchRequest>, I>>(object: I): CreateBatchRequest {
    const message = createBaseCreateBatchRequest();
    message.services = object.services?.map((e) => Service.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCreateBatchResponse(): CreateBatchResponse {
  return { services: [] };
}

export const CreateBatchResponse = {
  encode(message: CreateBatchResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.services) {
      Service.encode(v!, writer.uint32(10).fork()).ldelim();
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

          message.services.push(Service.decode(reader, reader.uint32()));
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
  async *encodeTransform(
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
  async *decodeTransform(
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
    return { services: Array.isArray(object?.services) ? object.services.map((e: any) => Service.fromJSON(e)) : [] };
  },

  toJSON(message: CreateBatchResponse): unknown {
    const obj: any = {};
    if (message.services?.length) {
      obj.services = message.services.map((e) => Service.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateBatchResponse>, I>>(base?: I): CreateBatchResponse {
    return CreateBatchResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateBatchResponse>, I>>(object: I): CreateBatchResponse {
    const message = createBaseCreateBatchResponse();
    message.services = object.services?.map((e) => Service.fromPartial(e)) || [];
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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
