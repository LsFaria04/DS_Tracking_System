process.env.MF_SKIP_SHARING_NEXT_INTERNAL_PLUGIN = 'true';
import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";
import { pluginModuleFederation } from "@module-federation/rsbuild-plugin";
import moduleFederationConfig from "./module-federation.config";

export default defineConfig({
  plugins: [pluginReact(), pluginModuleFederation(moduleFederationConfig)],
  server: {
    port: 5174,
    host: '0.0.0.0',
    headers: {
      'Access-Control-Allow-Origin': '*',
    },
  },
  source: {
    entry: {
      index: './src/index.tsx',
    },
  },
  html: {
    template: './index.html',
  },
});