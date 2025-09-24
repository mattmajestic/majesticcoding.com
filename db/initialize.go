package db

import "database/sql"

// InitializeDatabaseTables creates all necessary tables in the database
func InitializeDatabaseTables(dbConn *sql.DB) {
	CreateTables(dbConn)
	CreateMessagesTable(dbConn)
	CreateCheckinsTable(dbConn)
	CreateSpotifyTokensTable(dbConn)
	CreateTwitchTokensTable(dbConn)
	CreateTwitchMessagesTable(dbConn)
	CreateStatsHistoryTables(dbConn)
	CreateTwitchActivitiesTables(dbConn)
	CreateUsersTable(dbConn)
	CreateAuthSessionsTable(dbConn)

	// Vector tables for RAG
	CreateVectorTables(dbConn)
	CreateContextSummaryTable(dbConn)
}
