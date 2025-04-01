CREATE TABLE IF NOT EXISTS device_routines
(
    id varchar PRIMARY KEY NOT NULL,
    device_id varchar NOT NULL,
    status varchar NOT NULL,
    context varchar NOT NULL,
    area varchar NOT NULL,
    dispatched_at timestamp NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS device_routine_diagnostics 
(
    id varchar PRIMARY KEY NOT NULL,
    diagnostic varchar NOT NULL,
    routine_id varchar NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    CONSTRAINT fk_routine FOREIGN KEY (routine_id) REFERENCES device_routines (id) ON DELETE CASCADE ON UPDATE CASCADE
);

