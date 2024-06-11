DROP INDEX IF EXISTS idx_merchants_category_name;

DROP INDEX IF EXISTS idx_merchants_location;
DROP INDEX IF EXISTS idx_merchants_category;
DROP INDEX IF EXISTS idx_merchants_name;

DROP INDEX IF EXISTS idx_ll_to_earth;

DROP TABLE IF EXISTS merchants;

-- Drop extensions
DROP EXTENSION IF EXISTS "earthdistance";
DROP EXTENSION IF EXISTS "cube";
