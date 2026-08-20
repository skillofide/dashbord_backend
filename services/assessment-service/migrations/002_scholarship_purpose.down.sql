-- Rolling back has to deal with rows the old constraint cannot hold. Scholarship
-- papers become hiring drives rather than practice tests: 'hiring' is the value
-- that keeps them invite-only, so a rollback cannot accidentally throw a live
-- scholarship paper open to every logged-in student.
UPDATE assessments SET purpose = 'hiring' WHERE purpose = 'scholarship';

ALTER TABLE assessments DROP CONSTRAINT IF EXISTS assessments_purpose_check;
ALTER TABLE assessments ADD  CONSTRAINT assessments_purpose_check
    CHECK (purpose IN ('practice', 'hiring'));
