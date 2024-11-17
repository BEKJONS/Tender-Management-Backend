CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE  NOT NULL,
    password VARCHAR(255)        NOT NULL,
    role     VARCHAR(20)         NOT NULL CHECK (role IN ('client', 'contractor')),
    email    VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE tenders
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id   UUID REFERENCES users (id) ON DELETE CASCADE,
    title       VARCHAR(100) NOT NULL,
    description TEXT,
    deadline    TIMESTAMPTZ  NOT NULL,
    budget      NUMERIC(15, 2) CHECK (budget > 0),
    status      VARCHAR(20)      DEFAULT 'open' CHECK (status IN ('open', 'closed', 'awarded'))
);

CREATE TABLE bids
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id     UUID REFERENCES tenders (id) ON DELETE CASCADE,
    contractor_id UUID REFERENCES users (id) ON DELETE CASCADE,
    price         NUMERIC(15, 2) CHECK (price > 0),
    delivery_time INT CHECK (delivery_time > 0),
    comments      TEXT,
    status        VARCHAR(20)      DEFAULT 'pending'  CHECK (status IN ('pending', 'lost', 'awarded'))
);

CREATE TABLE notifications
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users (id) ON DELETE CASCADE,
    message     TEXT        NOT NULL,
    relation_id UUID,
    type        VARCHAR(20) NOT NULL,
    created_at  TIMESTAMPTZ      DEFAULT NOW()
);
