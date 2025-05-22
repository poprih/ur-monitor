-- Drop skcs table and its dependencies
DROP TABLE IF EXISTS skcs;

-- Revert units table to reference areas
ALTER TABLE units 
    DROP CONSTRAINT IF EXISTS fk_skc,
    DROP COLUMN IF EXISTS skc_id,
    ADD COLUMN area_id INTEGER,
    ADD CONSTRAINT fk_area
        FOREIGN KEY (area_id)
        REFERENCES areas(id)
        ON DELETE CASCADE;

-- Recreate index for area_id
CREATE INDEX IF NOT EXISTS idx_units_area_id ON units(area_id); 
