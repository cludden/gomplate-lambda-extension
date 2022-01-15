import { readFile } from 'fs/promises';

const config = JSON.parse(await readFile("/tmp/config.json"));

export async function handler(event, context) {
  return config;
}