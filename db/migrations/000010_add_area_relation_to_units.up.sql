ALTER TABLE units 
ADD COLUMN IF NOT EXISTS area_id INTEGER,
ADD CONSTRAINT fk_area
    FOREIGN KEY (area_id)
    REFERENCES areas(id)
    ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_units_area_id ON units(area_id); 
