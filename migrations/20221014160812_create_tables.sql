-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currency (
    id char(3) not null primary key,
    display VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS currency_rates (
    currency_id char(3) not null references currency(id),
    date timestamp not null,
    rate real,
    PRIMARY KEY (currency_id, date)
);

CREATE INDEX currency_rates_date on currency_rates(date);

CREATE TABLE IF NOT EXISTS user_states (
    user_id bigint not null primary key,
    currency_id char(3) not null references currency(id),
    status int not null default 0,
    input_buffer jsonb
);

CREATE TABLE IF NOT EXISTS categories (
    id serial not null primary key,
    name varchar(100) not null
);

CREATE TABLE IF NOT EXISTS expenses (
    id serial not null primary key,
    user_id bigint not null,
    category_id int not null references categories(id),
    currency_id char(3) not null references currency(id),
    amount int not null default 0,
    date timestamp not null
);

CREATE INDEX expenses_date ON expenses(date);

CREATE TABLE IF NOT EXISTS expense_limits (
    user_id bigint not null,
    month int not null check ( month > 0), check ( month < 13 ),
    value int not null,
    PRIMARY KEY (user_id, month)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE expense_limits;
DROP TABLE IF EXISTS expense_limits;

DROP INDEX IF EXISTS expenses_date;
TRUNCATE TABLE expenses;
DROP TABLE IF EXISTS expenses;

TRUNCATE TABLE categories;
DROP TABLE IF EXISTS categories;

TRUNCATE TABLE user_states;
DROP TABLE IF EXISTS user_states;

TRUNCATE TABLE currency_rates;
DROP INDEX IF EXISTS currency_rates_date;
DROP TABLE IF EXISTS currency_rates;

TRUNCATE TABLE currency;
DROP TABLE IF EXISTS currency;
-- +goose StatementEnd
