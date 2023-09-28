/* eslint-disable */

export const protobufPackage = "checker";

export const Action = { ACTION_UNSPECIFIED: 0, ACTION_CHECK: 1, ACTION_PUT: 2, ACTION_GET: 3 } as const;

export type Action = typeof Action[keyof typeof Action];

export function actionFromJSON(object: any): Action {
  switch (object) {
    case 0:
    case "ACTION_UNSPECIFIED":
      return Action.ACTION_UNSPECIFIED;
    case 1:
    case "ACTION_CHECK":
      return Action.ACTION_CHECK;
    case 2:
    case "ACTION_PUT":
      return Action.ACTION_PUT;
    case 3:
    case "ACTION_GET":
      return Action.ACTION_GET;
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum Action");
  }
}

export function actionToJSON(object: Action): string {
  switch (object) {
    case Action.ACTION_UNSPECIFIED:
      return "ACTION_UNSPECIFIED";
    case Action.ACTION_CHECK:
      return "ACTION_CHECK";
    case Action.ACTION_PUT:
      return "ACTION_PUT";
    case Action.ACTION_GET:
      return "ACTION_GET";
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum Action");
  }
}

export const Status = {
  STATUS_UNSPECIFIED: 0,
  STATUS_UP: 1,
  STATUS_MUMBLE: 2,
  STATUS_CORRUPT: 3,
  STATUS_DOWN: 4,
  STATUS_CHECK_FAILED: 5,
} as const;

export type Status = typeof Status[keyof typeof Status];

export function statusFromJSON(object: any): Status {
  switch (object) {
    case 0:
    case "STATUS_UNSPECIFIED":
      return Status.STATUS_UNSPECIFIED;
    case 1:
    case "STATUS_UP":
      return Status.STATUS_UP;
    case 2:
    case "STATUS_MUMBLE":
      return Status.STATUS_MUMBLE;
    case 3:
    case "STATUS_CORRUPT":
      return Status.STATUS_CORRUPT;
    case 4:
    case "STATUS_DOWN":
      return Status.STATUS_DOWN;
    case 5:
    case "STATUS_CHECK_FAILED":
      return Status.STATUS_CHECK_FAILED;
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum Status");
  }
}

export function statusToJSON(object: Status): string {
  switch (object) {
    case Status.STATUS_UNSPECIFIED:
      return "STATUS_UNSPECIFIED";
    case Status.STATUS_UP:
      return "STATUS_UP";
    case Status.STATUS_MUMBLE:
      return "STATUS_MUMBLE";
    case Status.STATUS_CORRUPT:
      return "STATUS_CORRUPT";
    case Status.STATUS_DOWN:
      return "STATUS_DOWN";
    case Status.STATUS_CHECK_FAILED:
      return "STATUS_CHECK_FAILED";
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum Status");
  }
}

export const Type = { TYPE_UNSPECIFIED: 0, TYPE_LEGACY: 1 } as const;

export type Type = typeof Type[keyof typeof Type];

export function typeFromJSON(object: any): Type {
  switch (object) {
    case 0:
    case "TYPE_UNSPECIFIED":
      return Type.TYPE_UNSPECIFIED;
    case 1:
    case "TYPE_LEGACY":
      return Type.TYPE_LEGACY;
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum Type");
  }
}

export function typeToJSON(object: Type): string {
  switch (object) {
    case Type.TYPE_UNSPECIFIED:
      return "TYPE_UNSPECIFIED";
    case Type.TYPE_LEGACY:
      return "TYPE_LEGACY";
    default:
      throw new tsProtoGlobalThis.Error("Unrecognized enum value " + object + " for enum Type");
  }
}

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
