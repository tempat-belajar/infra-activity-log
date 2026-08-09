ALTER TABLE activity_logs DROP CONSTRAINT IF EXISTS activity_logs_pic_check;
ALTER TABLE activity_logs ADD CONSTRAINT activity_logs_job_title_check
    CHECK (job_title IN ('DBA','DCO','NETWORK'));
