import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: '../contracts/openapi.bundled.yaml',
    output: {
      mode: 'tags-split',
      target: 'src/generated/api.ts',
      schemas: 'src/generated/models',
      client: 'react-query',
      mock: false,
      override: {
        mutator: {
          path: 'src/lib/axios-instance.ts',
          name: 'customInstance',
        },
      },
    },
  },
});