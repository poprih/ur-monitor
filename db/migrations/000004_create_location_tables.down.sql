DROP TRIGGER IF EXISTS update_areas_updated_at ON areas;
DROP TRIGGER IF EXISTS update_prefectures_updated_at ON prefectures;
DROP TRIGGER IF EXISTS update_regions_updated_at ON regions;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS areas;
DROP TABLE IF EXISTS prefectures;
DROP TABLE IF EXISTS regions; 
