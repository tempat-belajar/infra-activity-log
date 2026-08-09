-- Master tables for dropdown data
CREATE TABLE IF NOT EXISTS master_job_titles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_pics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(7),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default data
INSERT INTO master_job_titles (name) VALUES 
    ('DBA'), ('DCO'), ('NETWORK')
ON CONFLICT (name) DO NOTHING;

INSERT INTO master_pics (name) VALUES 
    ('Zaqi'), ('Dwi'), ('Irwan'), ('Kristopelly'), ('Benny')
ON CONFLICT (name) DO NOTHING;

INSERT INTO master_statuses (name, color) VALUES 
    ('Open', '#856404'),
    ('Process', '#0c5460'),
    ('Hold', '#6c757d'),
    ('Done', '#155724')
ON CONFLICT (name) DO NOTHING;

INSERT INTO master_categories (name) VALUES 
    ('Change'), ('Daily'), ('Incident')
ON CONFLICT (name) DO NOTHING;
