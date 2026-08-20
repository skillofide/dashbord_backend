-- Scholarship papers are a third kind of assessment.
--
-- They behave like a hiring drive — invite-only, one attempt, results decide an
-- outcome — but they are not a company's drive, so folding them into 'hiring'
-- would leave every recruiter report and company join carrying rows that belong
-- to no company. A distinct purpose keeps that segmentation honest.
--
-- Postgres names an inline CHECK "<table>_<column>_check"; DROP ... IF EXISTS
-- keeps this migration safe on a database where the constraint was already
-- renamed or dropped by hand.
ALTER TABLE assessments DROP CONSTRAINT IF EXISTS assessments_purpose_check;
ALTER TABLE assessments ADD  CONSTRAINT assessments_purpose_check
    CHECK (purpose IN ('practice', 'hiring', 'scholarship'));
