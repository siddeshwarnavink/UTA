import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';

export default defineConfig({
  mode: 'development',
  html: {
    template: './public/index.html',
  },
  plugins: [pluginReact({
    fastRefresh: false
  })],
});
