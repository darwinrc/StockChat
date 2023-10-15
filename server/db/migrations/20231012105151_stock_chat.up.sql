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

INSERT INTO users (id, username, password) VALUES ('48ccb5c1-9a19-42cd-bd41-3ac5c8af1108', 'StockBot', '$2a$13$x4z0aJ39z6e6Bo88m70/6uHLH.Kzo0KssEtq9BLZurL1Pe.SS/zDC');
INSERT INTO users (username, password) VALUES ('UserOne', '$2a$13$x4z0aJ39z6e6Bo88m70/6uHLH.Kzo0KssEtq9BLZurL1Pe.SS/zDC');
INSERT INTO users (username, password) VALUES ('UserTwo', '$2a$13$x4z0aJ39z6e6Bo88m70/6uHLH.Kzo0KssEtq9BLZurL1Pe.SS/zDC');