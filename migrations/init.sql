CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    borrower_id VARCHAR(100) NOT NULL,
    principal_amount NUMERIC(12,2) NOT NULL,
    rate NUMERIC(5,2) NOT NULL,
    roi NUMERIC(5,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'proposed',
    agreement_letter_link TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE loan_approvals (
    id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES loans(id) ON DELETE CASCADE,
    picture_proof TEXT NOT NULL,
    employee_id VARCHAR(50) NOT NULL,
    approved_at TIMESTAMP NOT NULL
);

CREATE TABLE investments (
    id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES loans(id) ON DELETE CASCADE,
    investor_email VARCHAR(100) NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    invested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);