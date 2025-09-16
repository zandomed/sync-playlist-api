-- migrations/001_initial_schema.sql

-- Extensión para UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum types
CREATE TYPE migration_status AS ENUM ('pending', 'running', 'completed', 'failed', 'cancelled');
CREATE TYPE migration_track_status AS ENUM ('pending', 'processing', 'success', 'failed', 'not_found');

-- Tabla de usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de autenticación de servicios
CREATE TABLE service_auths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_name VARCHAR(50) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMPTZ,
    scope TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, service_name)
);

-- Tabla de playlists
CREATE TABLE playlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_name VARCHAR(50) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    public BOOLEAN DEFAULT false,
    track_count INTEGER DEFAULT 0,
    owner_id VARCHAR(255),
    owner_name VARCHAR(255),
    image_url TEXT,
    external_url TEXT,
    last_synced_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, service_name, external_id)
);

-- Tabla de tracks
CREATE TABLE tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    external_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    duration INTEGER, -- en milisegundos
    isrc VARCHAR(20), -- Código internacional de grabación
    position INTEGER DEFAULT 0,
    image_url TEXT,
    preview_url TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de migraciones
CREATE TABLE migrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    target_service VARCHAR(50) NOT NULL,
    target_playlist_id UUID REFERENCES playlists(id),
    status migration_status DEFAULT 'pending',
    total_tracks INTEGER DEFAULT 0,
    processed_tracks INTEGER DEFAULT 0,
    successful_tracks INTEGER DEFAULT 0,
    failed_tracks INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de migración de tracks individuales
CREATE TABLE migration_tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    migration_id UUID NOT NULL REFERENCES migrations(id) ON DELETE CASCADE,
    source_track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    target_external_id VARCHAR(255),
    status migration_track_status DEFAULT 'pending',
    error_message TEXT,
    attempt_count INTEGER DEFAULT 0,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Índices para performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_service_auths_user_service ON service_auths(user_id, service_name);
CREATE INDEX idx_service_auths_expires ON service_auths(expires_at);

CREATE INDEX idx_playlists_user_service ON playlists(user_id, service_name);
CREATE INDEX idx_playlists_external ON playlists(service_name, external_id);
CREATE INDEX idx_playlists_sync ON playlists(last_synced_at);

CREATE INDEX idx_tracks_playlist ON tracks(playlist_id);
CREATE INDEX idx_tracks_external ON tracks(external_id);
CREATE INDEX idx_tracks_isrc ON tracks(isrc) WHERE isrc IS NOT NULL;
CREATE INDEX idx_tracks_search ON tracks USING GIN (to_tsvector('english', name || ' ' || artist));

CREATE INDEX idx_migrations_user ON migrations(user_id);
CREATE INDEX idx_migrations_status ON migrations(status);
CREATE INDEX idx_migrations_created ON migrations(created_at DESC);

CREATE INDEX idx_migration_tracks_migration ON migration_tracks(migration_id);
CREATE INDEX idx_migration_tracks_status ON migration_tracks(status);
CREATE INDEX idx_migration_tracks_source ON migration_tracks(source_track_id);

-- Triggers para updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_service_auths_updated_at BEFORE UPDATE ON service_auths
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_playlists_updated_at BEFORE UPDATE ON playlists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_migrations_updated_at BEFORE UPDATE ON migrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();