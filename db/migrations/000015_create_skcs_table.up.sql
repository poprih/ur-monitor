-- Create skcs table
CREATE TABLE skcs (
    id SERIAL PRIMARY KEY,
    area_id INTEGER NOT NULL REFERENCES areas(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    ur_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add trigger for updated_at
CREATE TRIGGER update_skcs_updated_at
    BEFORE UPDATE ON skcs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create index for area_id
CREATE INDEX idx_skcs_area_id ON skcs(area_id);

-- Modify units table to reference skcs instead of areas
ALTER TABLE units 
    DROP CONSTRAINT IF EXISTS fk_area,
    DROP COLUMN IF EXISTS area_id,
    ADD COLUMN skc_id INTEGER,
    ADD CONSTRAINT fk_skc
        FOREIGN KEY (skc_id)
        REFERENCES skcs(id)
        ON DELETE CASCADE;

-- Create index for skc_id
CREATE INDEX idx_units_skc_id ON units(skc_id); 
