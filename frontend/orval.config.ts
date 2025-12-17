import { defineConfig } from 'orval';

export default defineConfig({
  public: {
    input: '../backend/openapi/openapi-public.yaml',
    output: {
      target: './src/api/public',
      schemas: './src/api/public/model',
      client: 'react-query',
      mode: 'tags-split',
      override: {
        mutator: {
          path: './src/api/axios-instance.ts',
          name: 'customInstance',
        },
      },
    },
  },
});
