--
-- Add blocked_consumers table
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    CREATE TABLE blocked_consumers (
        organizer_id UUID NOT NULL,
        consumer_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        PRIMARY KEY (organizer_id, consumer_id),
        CONSTRAINT fk_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id),
        CONSTRAINT fk_consumer FOREIGN KEY (consumer_id) REFERENCES consumers(id)
    );
    
    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;