import type { Config } from 'jest';

const config: Config = {
  testEnvironment: 'jsdom',
  testRegex: '.*\\.spec\\.(ts|tsx|js|jsx)$',
  moduleNameMapper: {
    '\\.css$': '<rootDir>/__mocks__/styleMock.ts',
  },
  transform: {
    '^.+\\.[jt]sx?$': [
      '@swc/jest',
      {
        jsc: { transform: { react: { runtime: 'automatic' } } },
        module: { type: 'commonjs', noInterop: true },
        minify: false,
      },
    ],
  },
  setupFilesAfterEnv: ['./setup-tests.ts'],
};

export default config;
