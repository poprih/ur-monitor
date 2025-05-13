import { neon } from "@neondatabase/serverless";

const DATABASE_URL = import.meta.env.DATABASE_URL;

export const getRegions = async () => {
  const sql = neon(DATABASE_URL);
  const regions = await sql`SELECT * FROM regions`;
  return regions;
};

export const getPerfectures = async (regionId: string) => {
  const sql = neon(DATABASE_URL);
  const prefectures =
    await sql`SELECT * FROM prefectures WHERE region_id = ${regionId}`;
  return prefectures;
};
