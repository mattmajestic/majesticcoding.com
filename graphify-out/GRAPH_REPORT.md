# Graph Report - .  (2026-05-09)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 584 nodes · 926 edges · 54 communities (44 shown, 10 thin omitted)
- Extraction: 79% EXTRACTED · 21% INFERRED · 0% AMBIGUOUS · INFERRED: 197 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9714cec7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]

## God Nodes (most connected - your core abstractions)
1. `AIChatInterface` - 54 edges
2. `GetDB()` - 27 edges
3. `append()` - 26 edges
4. `SetupRoutes()` - 17 edges
5. `RedisSetJSON()` - 16 edges
6. `RedisGetJSON()` - 14 edges
7. `InitializeDatabaseTables()` - 14 edges
8. `CreateBronzeSchema()` - 13 edges
9. `main()` - 12 edges
10. `TranscribeAudio()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `getYouTubeStats()` --calls--> `InsertYouTubeStats()`  [INFERRED]
  api/handlers/stats.go → db/queries.go
- `StoreWebsiteContext()` --calls--> `GetDB()`  [INFERRED]
  api/services/embeddings.go → db/connect_db.go
- `StoreLatestSocialStatsContextFromDB()` --calls--> `GetDB()`  [INFERRED]
  api/services/embeddings.go → db/connect_db.go
- `StoreLatestSocialStatsContextFromDB()` --calls--> `GetLatestYouTubeStats()`  [INFERRED]
  api/services/embeddings.go → db/queries.go
- `StoreLatestSocialStatsContextFromDB()` --calls--> `GetLatestGitHubStats()`  [INFERRED]
  api/services/embeddings.go → db/queries.go

## Communities (54 total, 10 thin omitted)

### Community 1 - "Community 1"
Cohesion: 0.09
Nodes (21): Connect(), SaveTwitchToken(), UpdateTwitchToken(), SavedTwitchToken, TwitchUserLookupHandler(), exchangeTwitchCodeForToken(), GetTwitchUserToken(), InitTwitchClient() (+13 more)

### Community 2 - "Community 2"
Cohesion: 0.09
Nodes (22): append(), GetRecentMessages(), GetRecentTwitchMessages(), GetRecentTwitchUsersFromMessages(), InsertTwitchMessage(), StartMessageCleanup(), verifyClerkToken(), sendFinalTranscript() (+14 more)

### Community 3 - "Community 3"
Cohesion: 0.1
Nodes (19): cacheToken(), loadScript(), updateAuthUI(), CertificateList(), ServeCertificationPDF(), ClerkAuthMiddleware(), clerkFrontendDomain(), InitClerk() (+11 more)

### Community 4 - "Community 4"
Cohesion: 0.08
Nodes (21): interval, logo, progressBar, chatToggle, chatWidget, eyeClosed, eyeOpen, hiddenMsg (+13 more)

### Community 5 - "Community 5"
Cohesion: 0.1
Nodes (21): MapHandler(), CustomLogger(), sanitizeURL(), AddressComponents, GeocodeResult, GeoPoint, MapPage, gAddrComponent (+13 more)

### Community 6 - "Community 6"
Cohesion: 0.12
Nodes (21): GetTwitchToken(), ClearStatsCache(), StreamHandler(), StreamStatusHandler(), Stream, NewStream(), TwitchBits, TwitchRaid (+13 more)

### Community 7 - "Community 7"
Cohesion: 0.13
Nodes (20): appDiv, handleAuthenticatedUser(), showAuthenticatedState(), showUserButton(), waitForClerk(), aiMsg, appendMessage(), chatForm (+12 more)

### Community 8 - "Community 8"
Cohesion: 0.13
Nodes (20): InsertChatMessage(), GetProviders(), PostLLM(), LLMRequest, LLMResponse, callGemini(), callHuggingFace(), GenerateAIResponse() (+12 more)

### Community 9 - "Community 9"
Cohesion: 0.18
Nodes (17): GetDB(), GetRecentCheckins(), InsertCheckin(), InsertGitHubStats(), InsertTwitchStats(), PostCheckinHandler(), GeocodeHandler(), GlobeWidgetHandler() (+9 more)

### Community 10 - "Community 10"
Cohesion: 0.18
Nodes (14): GitHubStats, GraphQLError, GraphQLRequest, LeetCodeStats, TwitchStatsGQL, UnifiedStats, YouTubeStats, FetchGitHubStats() (+6 more)

### Community 11 - "Community 11"
Cohesion: 0.15
Nodes (10): openSocial(), model, hideServiceInfo(), highlightService(), resetView(), setupEventListeners(), setupServiceHovers(), showServiceInfo() (+2 more)

### Community 12 - "Community 12"
Cohesion: 0.17
Nodes (12): GetLatestGitHubStats(), GetLatestLeetCodeStats(), GetLatestTwitchStats(), GetLatestYouTubeStats(), InsertLeetCodeStats(), InsertTwitchBits(), InsertTwitchFollower(), InsertTwitchRaid() (+4 more)

### Community 13 - "Community 13"
Cohesion: 0.28
Nodes (13): CreateAuthSessionsTable(), CreateBronzeSchema(), CreateCheckinsTable(), CreateMessagesTable(), CreateSpotifyTokensTable(), CreateStatsHistoryTables(), CreateTables(), CreateTwitchActivitiesTables() (+5 more)

### Community 14 - "Community 14"
Cohesion: 0.19
Nodes (10): GetSpotifyToken(), SaveSpotifyToken(), SavedToken, exchangeCodeForToken(), loadToken(), refreshSpotifyToken(), saveToken(), SpotifyCallback() (+2 more)

### Community 15 - "Community 15"
Cohesion: 0.2
Nodes (13): CloudRunCostHandler(), CloudRunAverages, FetchCloudRunAverages(), firstTierUSD(), regionMatch(), safeAvg(), usdFromMoney(), money (+5 more)

### Community 16 - "Community 16"
Cohesion: 0.23
Nodes (12): chatUsers, createMetricCard(), dot, formatBytes(), formatMetricName(), formatUptime(), getMetricIcon(), loadMetricsCards() (+4 more)

### Community 17 - "Community 17"
Cohesion: 0.24
Nodes (10): SeedContext(), seedGitHubRepos(), seedSitePages(), ConvertStatsToText(), GenerateEmbedding(), RetrieveRelevantContext(), StoreSocialStatsContext(), StoreWebsiteContext() (+2 more)

### Community 18 - "Community 18"
Cohesion: 0.24
Nodes (9): InfracostHandler(), graphQLPrices, graphQLProduct, graphQLResp, escape(), FetchInfracostPrices(), InfracostPrice, InfracostRequest (+1 more)

### Community 19 - "Community 19"
Cohesion: 0.18
Nodes (6): main(), Execute(), LoadEnv(), StartBroadcaster(), InitSpotifyClient(), LoadMessagesFromDB()

### Community 20 - "Community 20"
Cohesion: 0.2
Nodes (7): distance, isDark, moonIcon, savedTheme, steps, sunIcon, toggleBtn

### Community 21 - "Community 21"
Cohesion: 0.36
Nodes (8): createProgressBar(), elements, formatTime(), loadSpotifyData(), startProgress(), stopProgress(), updateDisplay(), updateProgress()

### Community 22 - "Community 22"
Cohesion: 0.36
Nodes (7): isLiveNow(), offline, showOffline(), startPlayer(), stopPlayer(), tick(), video

### Community 23 - "Community 23"
Cohesion: 0.22
Nodes (8): GetRecentTwitchBits(), GetRecentTwitchFollowers(), GetRecentTwitchRaids(), GetRecentTwitchSubs(), TwitchBitsHandler(), TwitchRaidsHandler(), TwitchSubsHandler(), TwitchFollower

### Community 24 - "Community 24"
Cohesion: 0.25
Nodes (7): GetLaLigaSchedule(), LaLigaMatch, LaLigaResult, LaLigaScheduleResponse, LaLigaScore, LaLigaTeam, FetchLaLigaSchedule()

### Community 25 - "Community 25"
Cohesion: 0.25
Nodes (7): GetPremierLeagueSchedule(), PLMatch, PLResult, PLScheduleResponse, PLScore, PLTeam, FetchPLSchedule()

### Community 26 - "Community 26"
Cohesion: 0.54
Nodes (6): fetchMatches(), matches, rotateMatches(), showCurrentMatch(), showError(), showNoMatches()

### Community 27 - "Community 27"
Cohesion: 0.32
Nodes (6): adjustModalSize(), buttons, closeModal(), handleOutsideClick(), populateModalContent(), provider

### Community 28 - "Community 28"
Cohesion: 0.33
Nodes (4): activeButton, initializeTab(), showTab(), updateActiveButton()

### Community 29 - "Community 29"
Cohesion: 0.5
Nodes (3): copyButtons, textArea, url

### Community 31 - "Community 31"
Cohesion: 0.67
Nodes (3): getLatestCommit(), GitHashHandler(), GitCommit

### Community 32 - "Community 32"
Cohesion: 0.5
Nodes (3): GetCheckins(), GetCheckinsHandler(), Checkin

### Community 33 - "Community 33"
Cohesion: 0.67
Nodes (3): SaveScenario(), Scenario, NewScenario()

### Community 39 - "Community 39"
Cohesion: 0.67
Nodes (3): Google Cloud Infrastructure, Google Product Management, Google TensorFlow

### Community 40 - "Community 40"
Cohesion: 0.67
Nodes (3): Bitcoin Logo, NYU Machine Learning in Finance, Penn Analytics

## Knowledge Gaps
- **123 isolated node(s):** `chatToggle`, `hiddenMsg`, `eyeOpen`, `eyeClosed`, `input` (+118 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `append()` connect `Community 2` to `Community 32`, `Community 4`, `Community 6`, `Community 7`, `Community 8`, `Community 38`, `Community 9`, `Community 12`, `Community 15`, `Community 17`, `Community 18`, `Community 23`?**
  _High betweenness centrality (0.247) - this node is a cross-community bridge._
- **Why does `GetDB()` connect `Community 9` to `Community 32`, `Community 1`, `Community 2`, `Community 8`, `Community 12`, `Community 14`, `Community 17`, `Community 19`, `Community 23`?**
  _High betweenness centrality (0.171) - this node is a cross-community bridge._
- **Why does `contains()` connect `Community 5` to `Community 10`, `Community 15`?**
  _High betweenness centrality (0.169) - this node is a cross-community bridge._
- **Are the 26 inferred relationships involving `GetDB()` (e.g. with `StoreWebsiteContext()` and `StoreLatestSocialStatsContextFromDB()`) actually correct?**
  _`GetDB()` has 26 INFERRED edges - model-reasoned connections that need verification._
- **Are the 24 inferred relationships involving `append()` (e.g. with `ConvertStatsToText()` and `StoreLatestSocialStatsContextFromDB()`) actually correct?**
  _`append()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **Are the 15 inferred relationships involving `SetupRoutes()` (e.g. with `ClerkAuthMiddleware()` and `GeocodeHandler()`) actually correct?**
  _`SetupRoutes()` has 15 INFERRED edges - model-reasoned connections that need verification._
- **Are the 14 inferred relationships involving `RedisSetJSON()` (e.g. with `saveToken()` and `loadToken()`) actually correct?**
  _`RedisSetJSON()` has 14 INFERRED edges - model-reasoned connections that need verification._