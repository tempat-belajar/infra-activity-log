CREATE TABLE IF NOT EXISTS activity_logs (
    id SERIAL PRIMARY KEY,
    tanggal DATE NOT NULL,
    job_title VARCHAR(100),
    pic VARCHAR(20) NOT NULL CHECK (pic IN ('DBA','DCO','NETWORK')),
    application VARCHAR(150),
    label TEXT NOT NULL,
    old_value_text TEXT,
    old_value_image_url TEXT,
    new_value_text TEXT,
    new_value_image_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'Open' CHECK (status IN ('Open','Process','Hold','Done')),
    category VARCHAR(20) NOT NULL CHECK (category IN ('Change','Daily','Incident')),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_logs_tanggal ON activity_logs(tanggal DESC);
CREATE INDEX IF NOT EXISTS idx_logs_pic ON activity_logs(pic);
CREATE INDEX IF NOT EXISTS idx_logs_status ON activity_logs(status);
