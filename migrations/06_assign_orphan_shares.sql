-- The previous migration added a user_id column to the shares table
-- This migration will assign all shares that have no user_id to the first user in the database

-- Find the first user in the database
WITH first_user AS (
    SELECT id
    FROM users
    ORDER BY id
    LIMIT 1
)

-- Update shares with no user_id to the first user
UPDATE shares
SET user_id = (SELECT id FROM first_user)
WHERE user_id IS NULL;
