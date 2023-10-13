CREATE TABLE users
(
    id       uuid primary key default gen_random_uuid(),
    username text not null,
    password text not null
);

CREATE TABLE posts
(
    id        uuid primary key default gen_random_uuid(),
    user_id   uuid references users(id),
    message   text not null,
    timestamp timestamp not null default now()
);

INSERT INTO users (id, username, password) VALUES ('48ccb5c1-9a19-42cd-bd41-3ac5c8af1108', 'StockBot', '12345');
INSERT INTO users (username, password) VALUES ('UserOne', '12345');
INSERT INTO users (username, password) VALUES ('UserTwo', '12345');
INSERT INTO users (username, password) VALUES ('UserThree', '12345');