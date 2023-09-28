/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "pinger";

export interface PingRequest {
}

export interface PingResponse {
}

function createBasePingRequest(): PingRequest {
    return {};
}

export const PingRequest = {
    encode(_: PingRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
        return writer;
    },

    decode(input: _m0.Reader | Uint8Array, length?: number): PingRequest {
        const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePingRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
            }
            if ((tag & 7) === 4 || tag === 0) {
                break;
            }
            reader.skipType(tag & 7);
        }
        return message;
    },

    // encodeTransform encodes a source of message objects.
    // Transform<PingRequest, Uint8Array>
    async* encodeTransform(
        source: AsyncIterable<PingRequest | PingRequest[]> | Iterable<PingRequest | PingRequest[]>,
    ): AsyncIterable<Uint8Array> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [PingRequest.encode(p).finish()];
                }
            } else {
                yield* [PingRequest.encode(pkt).finish()];
            }
        }
    },

    // decodeTransform decodes a source of encoded messages.
    // Transform<Uint8Array, PingRequest>
    async* decodeTransform(
        source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
    ): AsyncIterable<PingRequest> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [PingRequest.decode(p)];
                }
            } else {
                yield* [PingRequest.decode(pkt)];
            }
        }
    },

    fromJSON(_: any): PingRequest {
        return {};
    },

    toJSON(_: PingRequest): unknown {
        const obj: any = {};
        return obj;
    },

    create<I extends Exact<DeepPartial<PingRequest>, I>>(base?: I): PingRequest {
        return PingRequest.fromPartial(base ?? ({} as any));
    },
    fromPartial<I extends Exact<DeepPartial<PingRequest>, I>>(_: I): PingRequest {
        const message = createBasePingRequest();
        return message;
    },
};

function createBasePingResponse(): PingResponse {
    return {};
}

export const PingResponse = {
    encode(_: PingResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
        return writer;
    },

    decode(input: _m0.Reader | Uint8Array, length?: number): PingResponse {
        const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePingResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
            }
            if ((tag & 7) === 4 || tag === 0) {
                break;
            }
            reader.skipType(tag & 7);
        }
        return message;
    },

    // encodeTransform encodes a source of message objects.
    // Transform<PingResponse, Uint8Array>
    async* encodeTransform(
        source: AsyncIterable<PingResponse | PingResponse[]> | Iterable<PingResponse | PingResponse[]>,
    ): AsyncIterable<Uint8Array> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [PingResponse.encode(p).finish()];
                }
            } else {
                yield* [PingResponse.encode(pkt).finish()];
            }
        }
    },

    // decodeTransform decodes a source of encoded messages.
    // Transform<Uint8Array, PingResponse>
    async* decodeTransform(
        source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
    ): AsyncIterable<PingResponse> {
        for await (const pkt of source) {
            if (Array.isArray(pkt)) {
                for (const p of pkt) {
                    yield* [PingResponse.decode(p)];
                }
            } else {
                yield* [PingResponse.decode(pkt)];
            }
        }
    },

    fromJSON(_: any): PingResponse {
        return {};
    },

    toJSON(_: PingResponse): unknown {
        const obj: any = {};
        return obj;
    },

    create<I extends Exact<DeepPartial<PingResponse>, I>>(base?: I): PingResponse {
        return PingResponse.fromPartial(base ?? ({} as any));
    },
    fromPartial<I extends Exact<DeepPartial<PingResponse>, I>>(_: I): PingResponse {
        const message = createBasePingResponse();
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
