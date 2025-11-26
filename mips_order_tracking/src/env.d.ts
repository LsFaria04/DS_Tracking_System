/// <reference types="@rsbuild/core/client" />

declare namespace NodeJS {
  interface ProcessEnv {
    readonly PUBLIC_API_URL: string;
    readonly PUBLIC_TOMTOM_API_KEY?: string;
  }
}