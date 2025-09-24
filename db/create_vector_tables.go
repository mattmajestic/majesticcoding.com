package db

import (
	"database/sql"
	"fmt"
)

// CreateVectorExtension enables pg_vector extension
func CreateVectorExtension(db *sql.DB) error {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS vector;`)
	if err != nil {
		return fmt.Errorf("failed to create vector extension: %w", err)
	}
	return nil
}

// CreateVectorTables creates tables for storing embeddings and context data
func CreateVectorTables(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	// Enable vector extension first
	if err := CreateVectorExtension(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		-- General website context table with embeddings for RAG
		CREATE TABLE IF NOT EXISTS bronze.website_context (
			id SERIAL PRIMARY KEY,
			content_type VARCHAR(50) NOT NULL, -- 'about_me', 'project', 'blog_post', 'skill', 'achievement', etc.
			title VARCHAR(255) NOT NULL, -- Title/name of the content
			content_text TEXT NOT NULL, -- Human-readable description for RAG
			source_url VARCHAR(500), -- Optional URL reference
			metadata JSONB, -- Additional structured data
			embedding vector(768), -- Gemini embeddings are 768 dimensions
			is_active BOOLEAN DEFAULT TRUE, -- Can disable content without deleting
			priority INTEGER DEFAULT 1, -- Higher priority content ranks higher
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		-- AI conversation context for better responses
		CREATE TABLE IF NOT EXISTS bronze.ai_conversations (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255),
			session_id VARCHAR(255),
			user_message TEXT NOT NULL,
			ai_response TEXT NOT NULL,
			context_used JSONB, -- What context was retrieved
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		-- Indexes for vector similarity search
		CREATE INDEX IF NOT EXISTS idx_website_context_embedding ON bronze.website_context
			USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

		CREATE INDEX IF NOT EXISTS idx_website_context_type ON bronze.website_context(content_type);
		CREATE INDEX IF NOT EXISTS idx_website_context_active ON bronze.website_context(is_active);
		CREATE INDEX IF NOT EXISTS idx_website_context_priority ON bronze.website_context(priority);

		-- Unique constraint for content_type + title
		CREATE UNIQUE INDEX IF NOT EXISTS idx_website_context_unique ON bronze.website_context(content_type, title);
		CREATE INDEX IF NOT EXISTS idx_ai_conversations_user_id ON bronze.ai_conversations(user_id);
		CREATE INDEX IF NOT EXISTS idx_ai_conversations_session_id ON bronze.ai_conversations(session_id);
	`)
	return err
}

// CreateContextSummaryTable for storing high-level summaries
func CreateContextSummaryTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		-- Summary table for quick context retrieval
		CREATE TABLE IF NOT EXISTS bronze.context_summaries (
			id SERIAL PRIMARY KEY,
			summary_type VARCHAR(50) NOT NULL, -- 'daily_stats', 'weekly_trends', 'achievements'
			title VARCHAR(255) NOT NULL,
			summary TEXT NOT NULL,
			data JSONB NOT NULL,
			embedding vector(1536),
			date_range_start TIMESTAMPTZ,
			date_range_end TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_context_summaries_embedding ON bronze.context_summaries
			USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

		CREATE INDEX IF NOT EXISTS idx_context_summaries_type ON bronze.context_summaries(summary_type);
		CREATE INDEX IF NOT EXISTS idx_context_summaries_date_range ON bronze.context_summaries(date_range_start, date_range_end);
	`)
	return err
}
