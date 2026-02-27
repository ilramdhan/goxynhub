-- Reset admin password to Admin@123456
-- Run this in Supabase SQL editor if you need to reset the admin password
-- This uses PostgreSQL's crypt() function which produces bcrypt hashes
-- compatible with Go's bcrypt library

-- Make sure pgcrypto extension is enabled
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Reset the admin user password
UPDATE users 
SET 
    password_hash = crypt('Admin@123456', gen_salt('bf', 12)),
    failed_attempts = 0,
    locked_until = NULL,
    status = 'active'
WHERE email = 'admin@example.com';

-- Verify the update
SELECT id, email, role, status, 
       (password_hash IS NOT NULL) as has_password
FROM users 
WHERE email = 'admin@example.com';
