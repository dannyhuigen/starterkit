-- name: GetAllGoogleUsers :many
SELECT *
FROM google_users;

-- name: UpsertGoogleUser :exec
INSERT INTO google_users (google_id, email, verified_email, name, picture, locale)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (google_id)
    DO UPDATE
    SET email = EXCLUDED.email,
        verified_email = EXCLUDED.verified_email,
        name = EXCLUDED.name,
        picture = EXCLUDED.picture,
        locale = EXCLUDED.locale;