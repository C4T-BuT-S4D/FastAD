/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "receiver";

export interface SubmitFlagsRequest {
  flags: string[];
}

export interface FlagResponse {
  flag: string;
  verdict: FlagResponse_Verdict;
  message: string;
  flagPoints: number;
}

export enum FlagResponse_Verdict {
  VERDICT_UNSPECIFIED = 0,
  VERDICT_ACCEPTED = 1,
  VERDICT_OWN = 2,
  VERDICT_OLD = 3,
  VERDICT_INVALID = 4,
  UNRECOGNIZED = -1,
}

export function flagResponse_VerdictFromJSON(
    object: any
): FlagResponse_Verdict {
  switch (object) {
    case 0:
    case "VERDICT_UNSPECIFIED":
      return FlagResponse_Verdict.VERDICT_UNSPECIFIED;
    case 1:
    case "VERDICT_ACCEPTED":
      return FlagResponse_Verdict.VERDICT_ACCEPTED;
    case 2:
    case "VERDICT_OWN":
      return FlagResponse_Verdict.VERDICT_OWN;
    case 3:
    case "VERDICT_OLD":
      return FlagResponse_Verdict.VERDICT_OLD;
    case 4:
    case "VERDICT_INVALID":
      return FlagResponse_Verdict.VERDICT_INVALID;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FlagResponse_Verdict.UNRECOGNIZED;
  }
}

export function flagResponse_VerdictToJSON(
    object: FlagResponse_Verdict
): string {
  switch (object) {
    case FlagResponse_Verdict.VERDICT_UNSPECIFIED:
      return "VERDICT_UNSPECIFIED";
    case FlagResponse_Verdict.VERDICT_ACCEPTED:
      return "VERDICT_ACCEPTED";
    case FlagResponse_Verdict.VERDICT_OWN:
      return "VERDICT_OWN";
    case FlagResponse_Verdict.VERDICT_OLD:
      return "VERDICT_OLD";
    case FlagResponse_Verdict.VERDICT_INVALID:
      return "VERDICT_INVALID";
    case FlagResponse_Verdict.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SubmitFlagsResponse {
  responses: FlagResponse[];
}

function createBaseSubmitFlagsRequest(): SubmitFlagsRequest {
  return {flags: []};
}

export const SubmitFlagsRequest = {
  encode(
      message: SubmitFlagsRequest,
      writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.flags) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubmitFlagsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubmitFlagsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.flags.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SubmitFlagsRequest {
    return {
      flags: Array.isArray(object?.flags)
          ? object.flags.map((e: any) => String(e))
          : [],
    };
  },

  toJSON(message: SubmitFlagsRequest): unknown {
    const obj: any = {};
    if (message.flags) {
      obj.flags = message.flags.map((e) => e);
    } else {
      obj.flags = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SubmitFlagsRequest>, I>>(
      object: I
  ): SubmitFlagsRequest {
    const message = createBaseSubmitFlagsRequest();
    message.flags = object.flags?.map((e) => e) || [];
    return message;
  },
};

function createBaseFlagResponse(): FlagResponse {
  return {flag: "", verdict: 0, message: "", flagPoints: 0};
}

export const FlagResponse = {
  encode(
      message: FlagResponse,
      writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.flag !== "") {
      writer.uint32(10).string(message.flag);
    }
    if (message.verdict !== 0) {
      writer.uint32(16).int32(message.verdict);
    }
    if (message.message !== "") {
      writer.uint32(26).string(message.message);
    }
    if (message.flagPoints !== 0) {
      writer.uint32(33).double(message.flagPoints);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FlagResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFlagResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.flag = reader.string();
          break;
        case 2:
          message.verdict = reader.int32() as any;
          break;
        case 3:
          message.message = reader.string();
          break;
        case 4:
          message.flagPoints = reader.double();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FlagResponse {
    return {
      flag: isSet(object.flag) ? String(object.flag) : "",
      verdict: isSet(object.verdict)
          ? flagResponse_VerdictFromJSON(object.verdict)
          : 0,
      message: isSet(object.message) ? String(object.message) : "",
      flagPoints: isSet(object.flagPoints) ? Number(object.flagPoints) : 0,
    };
  },

  toJSON(message: FlagResponse): unknown {
    const obj: any = {};
    message.flag !== undefined && (obj.flag = message.flag);
    message.verdict !== undefined &&
    (obj.verdict = flagResponse_VerdictToJSON(message.verdict));
    message.message !== undefined && (obj.message = message.message);
    message.flagPoints !== undefined && (obj.flagPoints = message.flagPoints);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FlagResponse>, I>>(
      object: I
  ): FlagResponse {
    const message = createBaseFlagResponse();
    message.flag = object.flag ?? "";
    message.verdict = object.verdict ?? 0;
    message.message = object.message ?? "";
    message.flagPoints = object.flagPoints ?? 0;
    return message;
  },
};

function createBaseSubmitFlagsResponse(): SubmitFlagsResponse {
  return {responses: []};
}

export const SubmitFlagsResponse = {
  encode(
      message: SubmitFlagsResponse,
      writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    for (const v of message.responses) {
      FlagResponse.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubmitFlagsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubmitFlagsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.responses.push(FlagResponse.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SubmitFlagsResponse {
    return {
      responses: Array.isArray(object?.responses)
          ? object.responses.map((e: any) => FlagResponse.fromJSON(e))
          : [],
    };
  },

  toJSON(message: SubmitFlagsResponse): unknown {
    const obj: any = {};
    if (message.responses) {
      obj.responses = message.responses.map((e) =>
          e ? FlagResponse.toJSON(e) : undefined
      );
    } else {
      obj.responses = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SubmitFlagsResponse>, I>>(
      object: I
  ): SubmitFlagsResponse {
    const message = createBaseSubmitFlagsResponse();
    message.responses =
        object.responses?.map((e) => FlagResponse.fromPartial(e)) || [];
    return message;
  },
};

type Builtin =
    | Date
    | Function
    | Uint8Array
    | string
    | number
    | boolean
    | undefined;

export type DeepPartial<T> = T extends Builtin
    ? T
    : T extends Array<infer U>
        ? Array<DeepPartial<U>>
        : T extends ReadonlyArray<infer U>
            ? ReadonlyArray<DeepPartial<U>>
            : T extends {}
                ? { [K in keyof T]?: DeepPartial<T[K]> }
                : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin
    ? P
    : P & { [K in keyof P]: Exact<P[K], I[K]> } & Record<Exclude<keyof I, KeysOfUnion<P>>,
    never>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
