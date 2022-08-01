CREATE TABLE sessions
(
    session_id   UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    expired_at     TIMESTAMP with time zone,
    id       UUID ,
    FOREIGN KEY (id)
        references users(id),
    archived_at timestamp  with time zone DEFAUlT NULL
);