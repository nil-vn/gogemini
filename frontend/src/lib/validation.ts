import type { ModuleKey, ModuleRecord, Settings } from './types';

const CURRENCY_RE = /^[A-Z]{3}$/;

export function validateSettings(input: Settings): string[] {
  const errors: string[] = [];
  if (!CURRENCY_RE.test(input.currency.trim())) errors.push('validationCurrency');
  if (input.theme !== 'light' && input.theme !== 'dark') errors.push('validationTheme');
  if (input.language !== 'en' && input.language !== 'vi') errors.push('validationLanguage');
  return errors;
}

export function validateRecord(module: ModuleKey, input: ModuleRecord): string[] {
  const errors: string[] = [];
  const name = String(input.name ?? '').trim();
  if ((module === 'users' || module === 'cars' || module === 'customers') && !name) errors.push('validationNameRequired');
  if (module === 'transactions' && !String(input.id ?? '').trim()) errors.push('validationTransactionIdRequired');
  return errors;
}

export function validateUploadFile(file: File | undefined): string[] {
  const errors: string[] = [];
  if (!file) return ['validationFileRequired'];
  if (!file.type.startsWith('image/')) errors.push('validationFileType');
  if (file.size > 5 * 1024 * 1024) errors.push('validationFileSize');
  return errors;
}
