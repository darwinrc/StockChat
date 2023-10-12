CREATE TABLE users
(
    id       uuid primary key default uuid_generate_v4(),
    username text not null,
    password text not null
);

CREATE TABLE posts
(
    id        uuid primary key default uuid_generate_v4(),
    user_id   uuid references users(id),
    message   text not null,
    timestamp timestamp not null default now()
);

INSERT INTO users (username, password)
VALUES ('StockBot', '12345');