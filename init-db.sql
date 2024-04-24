-- Create the databases
CREATE DATABASE tax_assessment;
CREATE DATABASE test_db;


-- Create a new user with a password
CREATE USER root WITH ENCRYPTED PASSWORD 'password';

-- Connect to the database tax_assessment to grant privileges
\c tax_assessment

-- Grant all privileges on database tax_assessment to the new user
GRANT ALL PRIVILEGES ON DATABASE tax_assessment TO root;
GRANT ALL PRIVILEGES ON SCHEMA public TO root;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO root;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO root;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO root;
-- Connect to the database test_db to grant privileges
\c test_db

-- Grant all privileges on database test_db to the new user
GRANT ALL PRIVILEGES ON DATABASE test_db TO root;
GRANT ALL PRIVILEGES ON SCHEMA public TO root;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO root;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO root;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO root;