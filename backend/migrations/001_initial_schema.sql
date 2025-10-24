-- Initial database schema for URL shortener application

-- Links table (using VARCHAR for ID instead of UUID to avoid extension dependency)
CREATE TABLE IF NOT EXISTS links (
    id VARCHAR(36) PRIMARY KEY,
    original_url TEXT NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(200),
    description TEXT,
    click_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Clicks table for analytics
CREATE TABLE IF NOT EXISTS clicks (
    id VARCHAR(36) PRIMARY KEY,
    link_id VARCHAR(36) REFERENCES links(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(10),
    device_type VARCHAR(20)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_links_slug ON links(slug);
CREATE INDEX IF NOT EXISTS idx_links_created_at ON links(created_at);
CREATE INDEX IF NOT EXISTS idx_clicks_link_id ON clicks(link_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at
CREATE TRIGGER update_links_updated_at 
    BEFORE UPDATE ON links 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();