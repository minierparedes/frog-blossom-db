DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        ALTER TABLE users
            DROP COLUMN IF EXISTS role;

        ALTER TABLE users
            ADD COLUMN IF NOT EXISTS role access NOT NULL;
    END IF;
END $$;
