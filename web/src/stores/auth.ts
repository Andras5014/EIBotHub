import { createPinia, defineStore } from 'pinia';

import { api } from '@/api';
import { ROLE_SUPER_ADMIN, resolvePermissions } from '@/constants/permissions';
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
    permissions: (state) => resolvePermissions(state.user ?? undefined),
    isAdmin: (state) => ['admin', ROLE_SUPER_ADMIN].includes(state.user?.role ?? ''),
    isSuperAdmin: (state) => state.user?.role === ROLE_SUPER_ADMIN,
    isOperator: (state) => state.user?.role === 'operator',
    isReviewer: (state) => state.user?.role === 'reviewer',
    hasBackofficeAccess: (state) => resolvePermissions(state.user ?? undefined).length > 0,
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
    hasRole(roles: string[]) {
      return this.isSuperAdmin || roles.includes(this.user?.role ?? '');
    },
    hasPermission(permission: string) {
      return this.isSuperAdmin || this.permissions.includes(permission);
    },
    hasAnyPermission(permissions: string[]) {
      return this.isSuperAdmin || permissions.some((permission) => this.permissions.includes(permission));
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
