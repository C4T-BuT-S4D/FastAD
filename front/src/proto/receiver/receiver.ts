/* eslint-disable */
import _m0 from "protobufjs/minimal";

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

export const FlagResponse_Verdict = {
    VERDICT_UNSPECIFIED: 0,
    VERDICT_ACCEPTED: 1,
    VERDICT_OWN: 2,
    VERDICT_OLD: 3,
    VERDICT_INVALID: 4,
} as const;

export type FlagResponse_Verdict = typeof FlagResponse_Verdict[keyof typeof FlagResponse_Verdict];

export function flagResponse_VerdictFromJSON(object: any): FlagResponse_Verdict {
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
        default:
            throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum FlagResponse_Verdict");
    }
}

export function flagResponse_VerdictToJSON(object: FlagResponse_Verdict): string {
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
        default:
            throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum FlagResponse_Verdict");
    }
}

export interface SubmitFlagsResponse {
    responses: FlagResponse[];
}

function createBaseSubmitFlagsRequest(): SubmitFlagsRequest {
    return {flags: []};
}

export const SubmitFlagsRequest = {
    encode(message: SubmitFlagsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
        for (const v of message.flags) {
            writer.uint32(10).string(v!);
        }
        return writer;
    },

    decode(input: _m0.Reader | Uint8Array, length?: number): SubmitFlagsRequest {
        const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSubmitFlagsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    if (tag !== 10) {
                        break;
                    }

                    message.flags.push(reader.string());
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
    // Transform<SubmitFlagsRequest, Uint8Array>
    async* encodeTransform(
        source:
            | AsyncIterable<SubmitFlagsRequest | SubmitFlagsRequest[]>
            | Iterable<SubmitFlagsRequest | SubmitFlagsRequest[]>,
    ): AsyncIterable<Uint8Array> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [SubmitFlagsRequest.encode(p).finish()];
                }
            } else {
                yield* [SubmitFlagsRequest.encode(pkt).finish()];
            }
        }
    },

    // decodeTransform decodes a source of encoded messages.
    // Transform<Uint8Array, SubmitFlagsRequest>
    async* decodeTransform(
        source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
    ): AsyncIterable<SubmitFlagsRequest> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [SubmitFlagsRequest.decode(p)];
                }
            } else {
                yield* [SubmitFlagsRequest.decode(pkt)];
            }
        }
    },

    fromJSON(object: any): SubmitFlagsRequest {
        return {flags: Array.isArray(object?.flags) ? object.flags.map((e: any) => String(e)) : []};
    },

    toJSON(message: SubmitFlagsRequest): unknown {
        const obj: any = {};
        if (message.flags?.length) {
            obj.flags = message.flags;
        }
        return obj;
    },

    create<I extends Exact<DeepPartial<SubmitFlagsRequest>, I>>(base?: I): SubmitFlagsRequest {
        return SubmitFlagsRequest.fromPartial(base ?? ({} as any));
    },
    fromPartial<I extends Exact<DeepPartial<SubmitFlagsRequest>, I>>(object: I): SubmitFlagsRequest {
        const message = createBaseSubmitFlagsRequest();
        message.flags = object.flags?.map((e) => e) || [];
        return message;
    },
};

function createBaseFlagResponse(): FlagResponse {
    return {flag: "", verdict: 0, message: "", flagPoints: 0};
}

export const FlagResponse = {
    encode(message: FlagResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
        const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFlagResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    if (tag !== 10) {
                        break;
                    }

                    message.flag = reader.string();
                    continue;
                case 2:
                    if (tag !== 16) {
                        break;
                    }

                    message.verdict = reader.int32() as any;
                    continue;
                case 3:
                    if (tag !== 26) {
                        break;
                    }

                    message.message = reader.string();
                    continue;
                case 4:
                    if (tag !== 33) {
                        break;
                    }

                    message.flagPoints = reader.double();
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
    // Transform<FlagResponse, Uint8Array>
    async* encodeTransform(
        source: AsyncIterable<FlagResponse | FlagResponse[]> | Iterable<FlagResponse | FlagResponse[]>,
    ): AsyncIterable<Uint8Array> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [FlagResponse.encode(p).finish()];
                }
            } else {
                yield* [FlagResponse.encode(pkt).finish()];
            }
        }
    },

    // decodeTransform decodes a source of encoded messages.
    // Transform<Uint8Array, FlagResponse>
    async* decodeTransform(
        source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
    ): AsyncIterable<FlagResponse> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [FlagResponse.decode(p)];
                }
            } else {
                yield* [FlagResponse.decode(pkt)];
            }
        }
    },

    fromJSON(object: any): FlagResponse {
        return {
            flag: isSet(object.flag) ? String(object.flag) : "",
            verdict: isSet(object.verdict) ? flagResponse_VerdictFromJSON(object.verdict) : 0,
            message: isSet(object.message) ? String(object.message) : "",
            flagPoints: isSet(object.flagPoints) ? Number(object.flagPoints) : 0,
        };
    },

    toJSON(message: FlagResponse): unknown {
        const obj: any = {};
        if (message.flag !== "") {
            obj.flag = message.flag;
        }
        if (message.verdict !== 0) {
            obj.verdict = flagResponse_VerdictToJSON(message.verdict);
        }
        if (message.message !== "") {
            obj.message = message.message;
        }
        if (message.flagPoints !== 0) {
            obj.flagPoints = message.flagPoints;
        }
        return obj;
    },

    create<I extends Exact<DeepPartial<FlagResponse>, I>>(base?: I): FlagResponse {
        return FlagResponse.fromPartial(base ?? ({} as any));
    },
    fromPartial<I extends Exact<DeepPartial<FlagResponse>, I>>(object: I): FlagResponse {
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
    encode(message: SubmitFlagsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
        for (const v of message.responses) {
            FlagResponse.encode(v!, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },

    decode(input: _m0.Reader | Uint8Array, length?: number): SubmitFlagsResponse {
        const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSubmitFlagsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    if (tag !== 10) {
                        break;
                    }

                    message.responses.push(FlagResponse.decode(reader, reader.uint32()));
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
    // Transform<SubmitFlagsResponse, Uint8Array>
    async* encodeTransform(
        source:
            | AsyncIterable<SubmitFlagsResponse | SubmitFlagsResponse[]>
            | Iterable<SubmitFlagsResponse | SubmitFlagsResponse[]>,
    ): AsyncIterable<Uint8Array> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [SubmitFlagsResponse.encode(p).finish()];
                }
            } else {
                yield* [SubmitFlagsResponse.encode(pkt).finish()];
            }
        }
    },

    // decodeTransform decodes a source of encoded messages.
    // Transform<Uint8Array, SubmitFlagsResponse>
    async* decodeTransform(
        source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
    ): AsyncIterable<SubmitFlagsResponse> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [SubmitFlagsResponse.decode(p)];
                }
            } else {
                yield* [SubmitFlagsResponse.decode(pkt)];
            }
        }
    },

    fromJSON(object: any): SubmitFlagsResponse {
        return {
            responses: Array.isArray(object?.responses) ? object.responses.map((e: any) => FlagResponse.fromJSON(e)) : [],
        };
    },

    toJSON(message: SubmitFlagsResponse): unknown {
        const obj: any = {};
        if (message.responses?.length) {
            obj.responses = message.responses.map((e) => FlagResponse.toJSON(e));
        }
        return obj;
    },

    create<I extends Exact<DeepPartial<SubmitFlagsResponse>, I>>(base?: I): SubmitFlagsResponse {
        return SubmitFlagsResponse.fromPartial(base ?? ({} as any));
    },
    fromPartial<I extends Exact<DeepPartial<SubmitFlagsResponse>, I>>(object: I): SubmitFlagsResponse {
        const message = createBaseSubmitFlagsResponse();
        message.responses = object.responses?.map((e) => FlagResponse.fromPartial(e)) || [];
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

function isSet(value: any): boolean {
    return value !== null && value !== undefined;
}
