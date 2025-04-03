ALTER TABLE areas 
ADD CONSTRAINT areas_prefecture_id_ur_area_code_unique 
UNIQUE (prefecture_id, ur_area_code); 
