-- Create the enum types
CREATE TYPE cat_races AS ENUM (
    'Persian',
    'Maine Coon',
    'Siamese',
    'Ragdoll',
    'Bengal',
    'Sphynx',
    'British Shorthair',
    'Abyssinian',
    'Scottish Fold',
    'Birman'
);

CREATE TYPE sex AS ENUM (
    'male',
    'female'
);

-- Create the cats table
CREATE TABLE cats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    race cat_races NOT NULL,
    sex sex NOT NULL,
    ageInMonth INTEGER NOT NULL,
    description VARCHAR NOT NULL,
    imageUrls VARCHAR[] NOT NULL,
    ownerId UUID NOT NULL REFERENCES users(id),
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    isDeleted BOOLEAN DEFAULT FALSE,
    hasMatched BOOLEAN DEFAULT FALSE
);
