CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(32)  NOT NULL,
    password   VARCHAR(255) NOT NULL,
    balance    INT          NOT NULL DEFAULT 1000 CHECK (balance >= 0),
    created_at TIMESTAMP    NOT NULL DEFAULT now(),
    CONSTRAINT username_unique UNIQUE (username)
);
CREATE INDEX idx_users_username ON users (username);


CREATE TABLE items
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(50) NOT NULL,
    price INT         NOT NULL CHECK (price > 0),
    CONSTRAINT name_unique UNIQUE (name)
);
CREATE INDEX idx_items_name ON items (name);

CREATE TABLE transactions
(
    id         SERIAL PRIMARY KEY,
    from_user  INT REFERENCES users (id),
    to_user    INT REFERENCES users (id),
    amount     INT       NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX idx_transactions_from_user ON transactions (from_user);
CREATE INDEX idx_transactions_to_user ON transactions (to_user);


CREATE TABLE inventory
(
    id         SERIAL PRIMARY KEY,
    user_id    INT       NOT NULL REFERENCES users (id),
    item_id    INT       NOT NULL REFERENCES items (id),
    quantity   INT       NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT unique_user_item UNIQUE (user_id, item_id)
);
CREATE INDEX idx_inventory_user_id ON inventory (user_id);
CREATE INDEX idx_inventory_item_id ON inventory (item_id);

INSERT INTO items (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);

