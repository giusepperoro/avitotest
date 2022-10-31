create extension if not exists pg_stat_statements;

-- Creation of accounts table
CREATE TABLE IF NOT EXISTS accounts (
    client_id BIGINT GENERATED ALWAYS AS IDENTITY
    (START WITH 1234 INCREMENT BY 1123),
    balance int NOT NULL,
    PRIMARY KEY(client_id)
);

-- Creation of orders table
CREATE TABLE IF NOT EXISTS orders (
    client_id BIGINT NOT NULL,
    order_id BIGINT UNIQUE NOT NULL,
    service_id int ,
    amount int,
    create_at timestamp NOT NULL,
    done_at timestamp,
    status int,
    PRIMARY KEY(order_id),
    CONSTRAINT fk_client
        FOREIGN KEY(client_id)
        REFERENCES accounts(client_id)
);
