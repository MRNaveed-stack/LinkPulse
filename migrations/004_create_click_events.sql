CREATE TABLE IF NOT EXISTS click_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  link_id UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE,
  is_bot BOOLEAN DEFAULT FALSE,
  user_agent TEXT,
  referrer TEXT,
  country VARCHAR(100),
  city VARCHAR(100),
  ip_address VARCHAR(45),
  clicked_at TIMESTAMP DEFAULT NOW()
);
