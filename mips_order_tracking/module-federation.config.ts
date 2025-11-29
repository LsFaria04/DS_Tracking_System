import { createModuleFederationConfig } from "@module-federation/rsbuild-plugin";

export default createModuleFederationConfig({
  name: 'mips_orders_provider',
  filename: 'remoteEntry.js',
  exposes: {
    './OrdersPage': './src/pages/HomePage.tsx',
    './OrderTracking': './src/pages/OrderTrackingPage.tsx',
  },
  shared: {
    react: { 
      singleton: true, 
      requiredVersion: '^18.0.0',
      eager: true,
      shareKey: 'react',
      shareScope: 'default'
    },
    'react-dom': { 
      singleton: true, 
      requiredVersion: '^18.0.0',
      eager: true,
      shareKey: 'react-dom',
      shareScope: 'default'
    },
    'react-router-dom': { 
      singleton: true,
      eager: true
    },
    'leaflet': {
      singleton: true,
      eager: true
    },
    'react-leaflet': {
      singleton: true,
      eager: true
    }
  },
});