import { NextResponse } from 'next/server';

export async function GET() {
  // Server-side routes can read regular env vars (not NEXT_PUBLIC_*)
  const apiUrl = process.env.API_URL || 'http://localhost:8080/';
  
  console.log('API Config - API_URL env var:', apiUrl);
  
  return NextResponse.json({
    apiUrl: apiUrl
  });
}
