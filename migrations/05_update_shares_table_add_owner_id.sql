PRAGMA foreign_keys = OFF;

-- We want to add a user_id column to the shares table
-- This will allow us to track the owner of the share
-- We will default the user_id to NULL
-- Due to sqlite limitations, we need to create a new table and copy the data over
-- Then we can drop the old table and rename the new one

-- Step 1: Create a new shares table with the foreign key constraint
CREATE TABLE shares_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path TEXT,
    expiration_date TEXT,
    long_id TEXT,
    num_files INTEGER,
    total_size INTEGER,
    files TEXT,
    user_id INTEGER, -- New foreign key column
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 2: Copy data from old table to new table (default user_id to NULL)
INSERT INTO shares_new (id, file_path, expiration_date, long_id, num_files, total_size, files)
SELECT id, file_path, expiration_date, long_id, num_files, total_size, files FROM shares;

-- Step 3: Drop old table
DROP TABLE shares;

-- Step 4: Rename new table to original name
ALTER TABLE shares_new RENAME TO shares;

PRAGMA foreign_keys = ON;