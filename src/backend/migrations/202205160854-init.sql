--
-- Initialize database
--
CREATE TABLE IF NOT EXISTS migrations (
    id VARCHAR PRIMARY KEY
);

SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE organizers (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name      TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP
    );

    CREATE TABLE products (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        organizer_id UUID NOT NULL,
        name TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,
        CONSTRAINT fk_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id)
    );

    CREATE TABLE bid_step_tables (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	    name      TEXT NOT NULL,
	    organizer_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,
        CONSTRAINT fk_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id)
    );

    CREATE TABLE consumers (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        nickname TEXT NOT NULL,
        form JSONB,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP
    );

    CREATE TABLE rooms (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        organizer_id UUID NOT NULL,
        name TEXT NOT NULL,
        address TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP
    );

    CREATE TABLE auctions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        room_id UUID NOT NULL,
        product_id UUID NOT NULL,
        bid_step_table_id UUID NOT NULL,
        min_amount DECIMAL(12, 2),
        started_at TIMESTAMP NOT NULL,
        finished_at TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,
        CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id),
        CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT fk_bid_step_table FOREIGN KEY (bid_step_table_id) REFERENCES bid_step_tables(id)
    );

    CREATE TABLE offers (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    	consumer_id UUID NOT NULL,
	    auction_id  UUID NOT NULL,
	    amount     DECIMAL(12, 2) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,
        CONSTRAINT fk_consumer FOREIGN KEY (consumer_id) REFERENCES consumers(id),
        CONSTRAINT fk_auction FOREIGN KEY (auction_id) REFERENCES auctions(id)
    );

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;