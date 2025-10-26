-- This migration removes the sample user.
DELETE FROM users WHERE email = 'test@example.com';