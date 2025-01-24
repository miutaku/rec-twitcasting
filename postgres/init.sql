CREATE DATABASE rec_twitcasting;

-- connection
\c rec_twitcasting;

-- create table
DROP TABLE IF EXISTS speakers;
CREATE TABLE speakers (
  id INTEGER PRIMARY KEY DEFAULT nextval('speakers_id_seq'),
  username VARCHAR(255) NOT NULL,
  recording_state BOOLEAN NOT NULL DEFAULT FALSE,
  created_date_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- create sequence
CREATE SEQUENCE speakers_id_seq START 1;

-- insert data (sample)
-- INSERT INTO rec_twitcasting.speakers (id, username) VALUES(nextval('speakers_id_seq'), 'twitcasting_username');
