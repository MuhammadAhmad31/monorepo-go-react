function getEnv<T extends string>(value: T | undefined, key: string): T {
  if (!value) {
    throw new Error(`Missing environment variable: ${key}`);
  }
  return value;
}

export const ENV = {
  BASE_URL: getEnv(import.meta.env.VITE_BASE_URL, 'VITE_BASE_URL'),
  ENVIRONMENT: getEnv(import.meta.env.VITE_ENVIRONMENT, 'VITE_ENVIRONMENT'),
} as const;
