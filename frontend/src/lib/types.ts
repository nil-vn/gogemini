export type ModuleKey = 'users' | 'cars' | 'customers' | 'transactions';

export type Settings = {
  currency: string;
  theme: 'dark' | 'light';
  language: 'vi' | 'en';
};

export type ModuleRecord = Record<string, unknown> & { id?: string | number };

export type DashboardMetrics = Record<string, number | string | unknown>;
