-- This migration adds a sample user for development and testing.
-- The password for this user is 'password123'.
-- The hash was generated using bcrypt with a cost of 10.
INSERT INTO users (username, email, password_hash)
VALUES ('testuser', 'test@example.com', '$2a$10$G.3gA.bL37pB/wPSHODLSOi9n2p5.L2x9S.wI3eNBh/r2iOnDXy/a');