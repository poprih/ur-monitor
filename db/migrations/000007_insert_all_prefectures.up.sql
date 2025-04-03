-- 北海道・東北
INSERT INTO prefectures (region_id, name, code, ur_code) 
SELECT id, '北海道', 'HOKKAIDO', '01' FROM regions WHERE code = 'HKTOHOKU'
UNION ALL
SELECT id, '宮城県', 'MIYAGI', '04' FROM regions WHERE code = 'HKTOHOKU';

-- 関東
INSERT INTO prefectures (region_id, name, code, ur_code) 
SELECT id, '東京都', 'TOKYO', '13' FROM regions WHERE code = 'KANTO'
UNION ALL
SELECT id, '神奈川県', 'KANAGAWA', '14' FROM regions WHERE code = 'KANTO'
UNION ALL
SELECT id, '千葉県', 'CHIBA', '12' FROM regions WHERE code = 'KANTO'
UNION ALL
SELECT id, '埼玉県', 'SAITAMA', '11' FROM regions WHERE code = 'KANTO'
UNION ALL
SELECT id, '茨城県', 'IBARAKI', '08' FROM regions WHERE code = 'KANTO';

-- 東海
INSERT INTO prefectures (region_id, name, code, ur_code) 
SELECT id, '愛知県', 'AICHI', '23' FROM regions WHERE code = 'TOKAI'
UNION ALL
SELECT id, '三重県', 'MIE', '24' FROM regions WHERE code = 'TOKAI'
UNION ALL
SELECT id, '岐阜県', 'GIFU', '21' FROM regions WHERE code = 'TOKAI'
UNION ALL
SELECT id, '静岡県', 'SHIZUOKA', '22' FROM regions WHERE code = 'TOKAI';

-- 関西
INSERT INTO prefectures (region_id, name, code, ur_code) 
SELECT id, '大阪府', 'OSAKA', '27' FROM regions WHERE code = 'KANSAI'
UNION ALL
SELECT id, '京都府', 'KYOTO', '26' FROM regions WHERE code = 'KANSAI'
UNION ALL
SELECT id, '兵庫県', 'HYOGO', '28' FROM regions WHERE code = 'KANSAI'
UNION ALL
SELECT id, '滋賀県', 'SHIGA', '25' FROM regions WHERE code = 'KANSAI'
UNION ALL
SELECT id, '奈良県', 'NARA', '29' FROM regions WHERE code = 'KANSAI'
UNION ALL
SELECT id, '和歌山県', 'WAKAYAMA', '30' FROM regions WHERE code = 'KANSAI';

-- 中国
INSERT INTO prefectures (region_id, name, code, ur_code) 
SELECT id, '岡山県', 'OKAYAMA', '33' FROM regions WHERE code = 'CHUGOKU'
UNION ALL
SELECT id, '広島県', 'HIROSHIMA', '34' FROM regions WHERE code = 'CHUGOKU'
UNION ALL
SELECT id, '山口県', 'YAMAGUCHI', '35' FROM regions WHERE code = 'CHUGOKU';

-- 九州
INSERT INTO prefectures (region_id, name, code, ur_code) 
SELECT id, '福岡県', 'FUKUOKA', '40' FROM regions WHERE code = 'KYUSHU'; 
