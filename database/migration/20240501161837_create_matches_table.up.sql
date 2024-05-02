CREATE TYPE match_statuses AS ENUM (
    'submitted',
    'cancelled',
    'approved',
    'rejected'
);

-- Create the matches table
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issuer_cat_id UUID NOT NULL REFERENCES cats(id),
    target_cat_id UUID NOT NULL REFERENCES cats(id),
    message VARCHAR NOT NULL,
    status match_statuses NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    isDeleted BOOLEAN DEFAULT FALSE,
    issuedBy UUID NOT NULL REFERENCES users(id)
);
