import { describe, expect, it } from 'vitest';
import { validateRecord, validateSettings, validateUploadFile } from '../validation';

describe('validation sync rules', () => {
  it('validates settings contract', () => {
    expect(validateSettings({ currency: 'usd', theme: 'light', language: 'en' })).toContain('validationCurrency');
    expect(validateSettings({ currency: 'USD', theme: 'dark', language: 'vi' })).toEqual([]);
  });

  it('validates record essentials by module', () => {
    expect(validateRecord('users', { id: '1' })).toContain('validationNameRequired');
    expect(validateRecord('transactions', { name: 'x' })).toContain('validationTransactionIdRequired');
  });

  it('validates upload constraints', () => {
    expect(validateUploadFile(undefined)).toContain('validationFileRequired');
    const bad = new File(['x'], 'x.txt', { type: 'text/plain' });
    expect(validateUploadFile(bad)).toContain('validationFileType');
  });
});
