INSERT INTO apps (name, auth_secret, refresh_secret)
VALUES ('test_app', 'test_app_auth_secret', 'test_app_refresh_secret')
ON CONFLICT DO NOTHING;
