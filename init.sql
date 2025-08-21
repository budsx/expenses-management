-- init.sql
-- Database initialization script for expenses management system

-- Create Users table with password column
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    role SMALLINT NOT NULL, -- 1=admin, 2=manager, 3=employee
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Expenses table
CREATE TABLE IF NOT EXISTS expenses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount_idr DECIMAL(15,2) NOT NULL,
    description TEXT NOT NULL,
    receipt_url VARCHAR(500),
    status SMALLINT NOT NULL DEFAULT 0, -- 0 Waiting for approval, 1 Approved, -1 Rejected
    auto_approved BOOLEAN DEFAULT FALSE,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- Create Approvals table
CREATE TABLE IF NOT EXISTS approvals (
    id BIGSERIAL PRIMARY KEY,
    expense_id BIGINT NOT NULL,
    approver_id BIGINT NOT NULL, -- user_id
    status SMALLINT NOT NULL, -- 1 Approved, -1 Rejected
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(expense_id, approver_id)
);

-- Create Expenses status log
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    expense_id BIGINT NOT NULL,
    new_status SMALLINT NOT NULL,
    status_before SMALLINT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_status ON expenses(status);
CREATE INDEX IF NOT EXISTS idx_expenses_submitted_at ON expenses(submitted_at);
CREATE INDEX IF NOT EXISTS idx_approvals_expense_id ON approvals(expense_id);
CREATE INDEX IF NOT EXISTS idx_approvals_approver_id ON approvals(approver_id);
CREATE INDEX IF NOT EXISTS idx_approvals_status ON approvals(status);
CREATE INDEX IF NOT EXISTS idx_audit_logs_expense_id ON audit_logs(expense_id);

-- Insert sample data with hashed passwords (bcrypt hash of "password123")
INSERT INTO users (email, name, role, password_hash) VALUES
    ('manager@company.com', 'Finance Manager', 2, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
    ('john.doe@company.com', 'John Doe', 3, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi')
ON CONFLICT (email) DO NOTHING;

-- Insert sample expenses
INSERT INTO expenses (user_id, amount_idr, description, receipt_url, status, auto_approved, submitted_at) VALUES
    (2, 150000.00, 'Lunch meeting with client', 'https://example.com/receipts/receipt1.jpg', 1, TRUE, NOW() - INTERVAL '2 days'),
    (2, 75000.00, 'Office supplies', 'https://example.com/receipts/receipt2.jpg', 1, TRUE, NOW() - INTERVAL '1 day'),
    (2, 200000.00, 'Taxi for business trip', 'https://example.com/receipts/receipt3.jpg', 0, TRUE, NOW())
ON CONFLICT DO NOTHING;

-- Insert sample approvals
INSERT INTO approvals (expense_id, approver_id, status, notes) VALUES
    (1, 1, 1, 'Approved lunch meeting expense'),
    (2, 1, 1, 'Approved for office supplies purchase')
ON CONFLICT DO NOTHING;