import axios, { AxiosError, type AxiosRequestConfig } from 'axios';

import type { Envelope } from '@/types/api';

const TOKEN_KEY = 'open-community-token';
const USER_KEY = 'open-community-user';

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
});

http.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (shouldResetSession(error)) {
      clearSessionAndRedirect();
    }
    return Promise.reject(error);
  },
);

export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  try {
    const response = await http.request<Envelope<T>>(config);
    return response.data.data;
  } catch (error) {
    throw toError(error);
  }
}

export async function sendWithoutData(config: AxiosRequestConfig): Promise<void> {
  try {
    await http.request(config);
  } catch (error) {
    throw toError(error);
  }
}

export { http };

function toError(error: unknown): Error {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<Envelope<never>>;
    const message = axiosError.response?.data?.error?.message ?? axiosError.message ?? 'request failed';
    return new Error(message);
  }
  return error instanceof Error ? error : new Error('request failed');
}

function shouldResetSession(error: unknown): boolean {
  if (!axios.isAxiosError(error)) {
    return false;
  }
  if (error.response?.status !== 401) {
    return false;
  }
  return Boolean(window.localStorage.getItem(TOKEN_KEY));
}

function clearSessionAndRedirect() {
  window.localStorage.removeItem(TOKEN_KEY);
  window.localStorage.removeItem(USER_KEY);

  if (window.location.pathname === '/login') {
    return;
  }

  const redirect = `${window.location.pathname}${window.location.search}${window.location.hash}`;
  window.location.replace(`/login?redirect=${encodeURIComponent(redirect)}`);
}
