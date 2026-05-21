import { describe, expect, it } from 'vitest';
import { isDebugJsonEditorEnabled } from '../debugAccess';

describe('isDebugJsonEditorEnabled', () => {
  it('returns false in production even when debug flag is true', () => {
    expect(isDebugJsonEditorEnabled({ debugJsonEditorFlag: 'true', isProd: true })).toBe(false);
  });

  it('returns true in non-production when debug flag is true', () => {
    expect(isDebugJsonEditorEnabled({ debugJsonEditorFlag: 'true', isProd: false })).toBe(true);
  });

  it('returns false when debug flag is not true', () => {
    expect(isDebugJsonEditorEnabled({ debugJsonEditorFlag: 'false', isProd: false })).toBe(false);
  });
});
