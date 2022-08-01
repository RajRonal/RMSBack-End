CREATE TABLE   users
(
    id            UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    user_name          TEXT NOT NULL,
    email         TEXT NOT NULL ,
    username      TEXT NOT NULL UNIQUE,
    password       TEXT NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at   TIMESTAMP WITH TIME ZONE


);
CREATE TYPE user_role AS ENUM ('user', 'admin','sub-admin');
CREATE TABLE  location
(
       longitude  DOUBLE PRECISION NOT NULL ,
       latitude    DOUBLE PRECISION NOT NULL ,
        created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        archived_at   TIMESTAMP WITH TIME ZONE,
       user_id  UUID REFERENCES users(id)
);
CREATE TABLE roles
(
    user_id     UUID NOT NULL REFERENCES users (id),
    user_role        user_role               DEFAULT 'user'::user_role,
    username    text not null,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
CREATE TABLE restaurant
(
    restaurant_id            UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    restaurant_name          TEXT NOT NULL,
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at   TIMESTAMP WITH TIME ZONE,
    longitude  DOUBLE PRECISION NOT NULL ,
    latitude    DOUBLE PRECISION NOT NULL
);
CREATE TABLE dishes
(
    dish_id            uuid PRIMARY KEY      DEFAULT gen_random_uuid(),
    dish_name          text  NOT NULL,
    dish_price         FLOAT NOT NULL,
    restaurant_id uuid  NOT NULL REFERENCES restaurant (restaurant_id),
    created_by    uuid REFERENCES users (id),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    archived_at   TIMESTAMP WITH TIME ZONE DEFAULT NULL
);