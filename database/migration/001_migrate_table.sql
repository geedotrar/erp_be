
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    position_name VARCHAR(255) NOT NULL UNIQUE,
    position_code VARCHAR(255) NOT NULL UNIQUE,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
);
CREATE TABLE company (
    id SERIAL PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL UNIQUE,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    role VARCHAR(10) CHECK (role IN ('user', 'admin')) NOT NULL,
    status BOOLEAN,
    position_id INT NOT NULL,
    company_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (position_id) REFERENCES positions(id),
    FOREIGN KEY (company_id) REFERENCES company(id)
);

INSERT INTO positions (position_name, position_code) VALUES
    ('Manager', 'MGR'),
    ('Supervisor', 'SPV');

INSERT INTO company (company_name) VALUES
    ('Company A'),
    ('Company B');

INSERT INTO users (first_name, last_name, email, password, phone_number, role, status, position_id, company_id)
VALUES
    ('John', 'Doe', 'john.doe@example.com', 'password123', '123456789', 'user', true, 1, 1),
    ('Jane', 'Smith', 'jane.smith@example.com', 'password456', '987654321', 'admin', false, 2, 2);