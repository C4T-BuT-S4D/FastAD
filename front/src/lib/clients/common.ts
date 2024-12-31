import {
  AxiosError,
  AxiosInstance,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';

const onRequest = (
  config: InternalAxiosRequestConfig,
): InternalAxiosRequestConfig => {
  return config;
};

const onRequestError = (error: AxiosError): Promise<AxiosError> => {
  console.error(`[request error] [${JSON.stringify(error)}]`);
  return Promise.reject(error);
};

const onResponse = (response: AxiosResponse): AxiosResponse => {
  return response;
};

const onResponseError = (error: AxiosError): Promise<AxiosError> => {
  console.error(`[response error] [${JSON.stringify(error)}]`);
  return Promise.reject(error);
};

export function setupInterceptorsTo(
  axiosInstance: AxiosInstance,
): AxiosInstance {
  axiosInstance.interceptors.request.use(onRequest, onRequestError);
  axiosInstance.interceptors.response.use(onResponse, onResponseError);
  return axiosInstance;
}

interface ProtoMessage {
  // eslint-disable-next-line no-unused-vars
  fromJSON(data: any): any;
  // eslint-disable-next-line no-unused-vars
  decode(data: Uint8Array): any;
}

export function createProtoJSONTransform<T extends ProtoMessage>(
  messageType: T,
) {
  return (response: string) => {
    return messageType.fromJSON(JSON.parse(response));
  };
}

export function createProtoTransform<T extends ProtoMessage>(messageType: T) {
  return (response: ArrayBufferLike) => {
    return messageType.decode(new Uint8Array(response));
  };
}
