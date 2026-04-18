import { createPinia, defineStore } from 'pinia';

import { api } from '@/api';
import type { UserSummary } from '@/types/api';

export const pinia = createPinia();

const TOKEN_KEY = 'open-community-token';
const USER_KEY = 'open-community-user';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) ?? '',
    user: readUser(),
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    isAdmin: (state) => state.user?.role === 'admin',
  },
  actions: {
    async login(payload: { email: string; password: string }) {
      const result = await api.login(payload);
      this.setAuth(result.token, result.user);
      return result;
    },
    async register(payload: { username: string; email: string; password: string }) {
      const result = await api.register(payload);
      this.setAuth(result.token, result.user);
      return result;
    },
    async refreshMe() {
      if (!this.token) return;
      this.user = await api.me();
      localStorage.setItem(USER_KEY, JSON.stringify(this.user));
    },
    async logout() {
      if (this.token) {
        try {
          await api.logout();
        } catch {
          // Ignore logout API failures and clear local session anyway.
        }
      }
      this.token = '';
      this.user = null;
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
    },
    setAuth(token: string, user: UserSummary) {
      this.token = token;
      this.user = user;
      localStorage.setItem(TOKEN_KEY, token);
      localStorage.setItem(USER_KEY, JSON.stringify(user));
    },
  },
});

function readUser(): UserSummary | null {
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as UserSummary;
  } catch {
    return null;
  }
}
