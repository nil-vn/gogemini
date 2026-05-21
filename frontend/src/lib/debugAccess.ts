export type DebugAccessEnv = {
  debugJsonEditorFlag: string | boolean | undefined;
  isProd: boolean;
};

export function isDebugJsonEditorEnabled(env: DebugAccessEnv): boolean {
  const flag = String(env.debugJsonEditorFlag ?? 'false').toLowerCase() === 'true';
  return flag && !env.isProd;
}
