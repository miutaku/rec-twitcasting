DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'rec_twitcasting') THEN
        CREATE DATABASE rec_twitcasting;
    END IF;
END
$$;

-- connection
\c rec_twitcasting;

-- create sequence
CREATE SEQUENCE speakers_id_seq START 1;

-- create table
CREATE TABLE IF NOT EXISTS speakers (
    id INTEGER DEFAULT nextval('speakers_id_seq'),
    username VARCHAR(255) PRIMARY KEY NOT NULL,
    recording_state BOOLEAN NOT NULL DEFAULT FALSE,
    created_date_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- insert data (sample)
-- INSERT INTO speakers (username) VALUES('twitcasting_username');

-- delete data (sample)
-- DELETE FROM speakers WHERE username = 'twitcasting_username';

-- update data (sample)
-- UPDATE speakers SET recording_state = TRUE WHERE username = 'twitcasting_username';