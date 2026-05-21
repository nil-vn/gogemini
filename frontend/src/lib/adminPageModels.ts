import type { ModuleRecord } from './types';

export type UserForm = { username: string; email: string; password: string; confirm_password: string; role: string; status: string };
export type CarForm = Record<string, string>;
export type CustomerForm = Record<string, string>;
export type TransactionForm = Record<string, string>;
export type TransactionItem = { name: string; price: string };

export function toUserForm(input: Partial<ModuleRecord> = {}): UserForm {
  return {
    username: String(input.username ?? ''),
    email: String(input.email ?? ''),
    password: '',
    confirm_password: '',
    role: String(input.role ?? 'guest') || 'guest',
    status: String(input.status ?? 'Active') || 'Active'
  };
}

export function validateUserForm(userForm: UserForm, isEdit = false): string {
  if (!userForm.username.trim()) return 'Username is required';
  if (userForm.password || userForm.confirm_password) {
    if (userForm.password !== userForm.confirm_password) return 'Password confirmation does not match';
  } else if (!isEdit) {
    return 'Password is required for new user';
  }
  return '';
}

export function toCarForm(input: Partial<ModuleRecord> = {}): CarForm { return {
  name: String(input.name ?? ''), model: String(input.model ?? ''), year_of_manufacture: String(input.year_of_manufacture ?? ''),
  vin: String(input.vin ?? ''), imported_date: String(input.imported_date ?? ''), purchase_price: String(input.purchase_price ?? ''), inspection_from: String(input.inspection_from ?? ''),
  status: String(input.status ?? 'AVAILABLE') || 'AVAILABLE', car_situation: String(input.car_situation ?? 'NOT_REFURBISHED') || 'NOT_REFURBISHED', color: String(input.color ?? ''),
  branch: String(input.branch ?? ''), license_plate_no: String(input.license_plate_no ?? ''), traded_company: String(input.traded_company ?? ''), selling_price: String(input.selling_price ?? ''),
  inspection_to: String(input.inspection_to ?? ''), note: String(input.note ?? '')
}; }

export function toCustomerForm(input: Partial<ModuleRecord> = {}): CustomerForm { return {
  name: String(input.name ?? ''), gender: String(input.gender ?? 'unknown') || 'unknown', address: String(input.address ?? ''), phone: String(input.phone ?? ''),
  birth_day: String(input.birth_day ?? ''), facebook: String(input.facebook ?? ''), lead_source: String(input.lead_source ?? ''), status: String(input.status ?? ''), note: String(input.note ?? '')
}; }

export function toTransactionForm(input: Partial<ModuleRecord> = {}): TransactionForm { return {
  customer_id: String((input as any).customer_id ?? (input as any).customer?.id ?? ''),
  car_id: String((input as any).car_id ?? (Array.isArray((input as any).cars) && (input as any).cars[0]?.id) ?? ''),
  purchase_date: String(input.purchase_date ?? ''), selling_price: String(input.selling_price ?? ''), deposit_amount: String(input.deposit_amount ?? ''), status: String(input.status ?? ''), note: String(input.note ?? '')
}; }

export function toTransactionItems(input: Partial<ModuleRecord> = {}): TransactionItem[] {
  const raw = Array.isArray((input as any).items) ? (input as any).items : [];
  const mapped = raw.map((it: any) => ({ name: String(it?.name ?? ''), price: String(it?.price ?? '') }));
  return mapped.length ? mapped : [{ name: '', price: '' }];
}

export function buildTransactionSummary(records: ModuleRecord[]) {
  const txs = records ?? [];
  const toNum = (v: unknown) => Number(v ?? 0) || 0;
  const isPaid = (tx: any) => String(tx.status ?? '').toLowerCase().includes('paid') || String(tx.status ?? '').includes('Đã');
  const totalRevenue = txs.reduce((acc, tx: any) => acc + toNum(tx.total_amount || tx.selling_price), 0);
  const paidRevenue = txs.filter((tx: any) => isPaid(tx)).reduce((acc, tx: any) => acc + toNum(tx.total_amount || tx.selling_price), 0);
  const depositedAmount = txs.filter((tx: any) => !isPaid(tx)).reduce((acc, tx: any) => acc + toNum(tx.deposit_amount), 0);
  return { txs, totalRevenue, paidRevenue, depositedAmount };
}
