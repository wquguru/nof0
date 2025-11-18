# Persistence Model Diagram Review

> **Diagram:** 07b-persistence-model.puml
> **Reviewer:** Architecture Review Team
> **Date:** 2025-11-18
> **Status:** 🟡 In Progress

---

## 1. Metadata Check ✅

- [x] **Title**: "Ideal Persistence Model (Database Tables)" - Clear ✅
- [x] **Theme**: `!theme plain` - Consistent ✅
- [x] **Legend**: Color-coded table types (Config/State/Event/Snapshot) ✅
- [x] **Notes**: Detailed notes for each module ✅

---

## 2. Content Review 📋

### 2.1 Completeness

- [x] All required entities present?
  - ✅ Manager Module: trader_config, positions, trades
  - ✅ Executor/LLM Module: models, conversation_messages, decision_cycles
  - ✅ Exchange/Account Module: accounts, account_snapshots
  - ✅ Exchange/Market Module: symbols, market_metrics, klines, price_ticks

- [x] All relationships defined?
  - ✅ FK relationships marked
  - ✅ Cross-module relationships shown

**Missing Elements:** None identified

### 2.2 Accuracy

**Code Verification:**

| Table | SQL File | Status | Notes |
|-------|----------|--------|-------|
| trader_config | `schema/trader_config.sql` | 🔍 Need to verify | - |
| positions | `schema/positions.sql` | 🔍 Need to verify | - |
| trades | `schema/trades.sql` | 🔍 Need to verify | - |
| decision_cycles | `schema/decision_cycles.sql` | 🔍 Need to verify | - |
| symbols | `schema/symbols.sql` | 🔍 Need to verify | - |

**Action Required:** Verify against actual SQL schema files in codebase

### 2.3 Clarity

- [x] Easy to understand at first glance? ✅
  - Color coding helps distinguish table types
  - Logical grouping by module
- [x] Logical grouping of components? ✅
  - Packages separate concerns well
- [x] Text readable? ✅
  - No overlaps observed

**Clarity Issues:** None

---

## 3. Domain-Specific Checks 🎯

### For Persistence Diagrams

- [x] Primary keys marked? ✅ (marked with `<<PK>>`)
- [x] Foreign keys marked? ✅ (marked with `<<FK>>`)
- [x] Indexes documented? ✅ (in notes)
- [x] JSONB fields explained? ✅ (detail JSONB contents documented)
- [x] Table types color-coded? ✅
  - #E3F2FD: Configuration
  - #FFF3E0: State
  - #E8F5E9: Event Log
  - #FCE4EC: Snapshot

**Specific Checks:**

#### Event Log Pattern Validation
- [x] `positions` table: Has `event_at`, append-only design ✅
- [x] `trades` table: Has `event_at`, immutable rows ✅
- [x] `decision_cycles` table: Has `executed_at`, event sourcing ✅
- [x] `klines` table: Has `open_time`, time-series data ✅

#### JSONB Usage Validation
- [x] `trader_config.detail`: Contains all nested config ✅
- [x] `positions.detail`: Contains position-specific data ✅
- [x] `trades.detail`: Contains execution details ✅
- [x] `decision_cycles.detail`: Contains cycle metadata ✅

#### Index Strategy
```
Documented Indexes:
✅ idx_trader_config_providers ON (exchange_provider, market_provider)
✅ idx_positions_trader_status ON (trader_id, status)
✅ idx_positions_symbol_open ON (symbol) WHERE status='open'
✅ idx_trades_trader_close_ts ON (trader_id, event_at DESC)
✅ idx_decision_cycles_trader_executed_at_desc ON (trader_id, executed_at DESC)
✅ idx_market_metrics_symbol_event_at_desc ON (symbol_id, event_at DESC)
✅ idx_klines_symbol_interval_open_time_desc ON (symbol_id, interval, open_time DESC)
✅ idx_price_ticks_symbol_event_at_desc ON (symbol, event_at DESC)
```

**Missing Indexes to Consider:**
- [ ] `trader_config(id)` - Already PK, auto-indexed
- [ ] `conversation_messages(conversation_id)` - For querying by conversation
- [ ] `symbols(exchange_provider, is_delisted)` - For filtering active symbols

---

## 4. Code Alignment Check 🔗

**Files to verify against:**
```bash
# SQL Schema files (if exist)
schema/migrations/*.sql

# Model definitions
pkg/repo/model/*.go
pkg/repo/query/*.sql.go

# Config structs
pkg/manager/config.go
pkg/executor/config.go
pkg/llm/config.go
```

**Verification Status:**
- 🔍 **Action Required**: Map each table to corresponding Go struct
- 🔍 **Action Required**: Verify JSONB fields match Go struct tags

**Example Mapping to Verify:**
```go
// trader_config table → pkg/manager.TraderConfig struct
type TraderConfig struct {
    ID                string        `db:"id" json:"id"`
    ExchangeProvider  string        `db:"exchange_provider" json:"exchange_provider"`
    MarketProvider    string        `db:"market_provider" json:"market_provider"`
    Detail            TraderDetail  `db:"detail" json:"detail"` // JSONB
    DetailChecksum    int           `db:"detail_checksum" json:"detail_checksum"`
    CreatedAt         time.Time     `db:"created_at" json:"created_at"`
    UpdatedAt         time.Time     `db:"updated_at" json:"updated_at"`
}
```

---

## 5. Issues Found 🐛

### Critical Issues (Must Fix)
_None identified_

### Medium Issues (Should Fix)

1. **Missing partial index on `klines` for latest data**
   - **Location**: `klines` table indexes
   - **Impact**: Query performance for "get latest N klines"
   - **Suggested Fix**: Add `idx_klines_symbol_interval_latest ON (symbol_id, interval, open_time DESC) WHERE open_time > NOW() - INTERVAL '7 days'`

2. **No composite index for common queries**
   - **Location**: `decision_cycles` table
   - **Impact**: Queries filtering by `trader_id` + `error_message IS NULL` may be slow
   - **Suggested Fix**: Consider `idx_decision_cycles_trader_success ON (trader_id, executed_at DESC) WHERE error_message IS NULL`

### Minor Issues (Nice to Have)

1. **UNIQUE constraint documentation**
   - Some tables have UNIQUE constraints in notes but not explicitly marked
   - Consider adding `<<UK>>` stereotype

2. **Partition strategy not mentioned**
   - For large time-series tables (`klines`, `price_ticks`), consider documenting partition strategy
   - e.g., "Partition by month on `open_time`"

---

## 6. Cross-Diagram Consistency Check 🔗

**Entities mentioned in other diagrams:**

| Entity | This Diagram | 07a-domain-model.puml | Status |
|--------|-------------|----------------------|--------|
| Trader | ❌ Not a table | ✅ Aggregate Root | ✅ Correct (runtime only) |
| TraderConfig | ✅ trader_config | ✅ Entity | ✅ Match |
| Position | ✅ positions | ✅ Entity | ✅ Match |
| Decision | ❌ Not persisted | ✅ Value Object | ✅ Correct (ephemeral) |
| DecisionCycle | ✅ decision_cycles | ✅ Entity | ✅ Match |
| Symbol | ✅ symbols | ✅ Value Object | ⚠️ Mismatch (VO in domain but table here) |

**Analysis of Symbol discrepancy:**
- In domain model: `Symbol` is a Value Object (immutable, no identity)
- In persistence: `symbols` is a table with PK (has identity)
- **Resolution**: This is acceptable - Value Objects can be persisted as reference data
- **Suggestion**: Add note in 07a explaining this design choice

---

## 7. Recommendations 💡

### Improvements

- [x] **Add data retention policy notes**
  - For event logs: How long to keep `decision_cycles`, `trades`?
  - For snapshots: Aggregation strategy for `account_snapshots`, `market_metrics`?

- [x] **Document JSONB schema versions**
  - Add `detail_schema_version` field to tables with JSONB
  - Enables forward-compatible schema evolution

- [x] **Add soft-delete pattern**
  - For `trader_config`: Add `deleted_at` for soft delete instead of hard delete
  - Preserves audit trail

### Questions for Team

1. **Event Log Retention**: What's the retention policy for `decision_cycles` and `trades`?
   - Suggestion: Archive to cold storage after 90 days, keep summary stats

2. **JSONB vs Columns**: Should frequently-queried JSONB fields be promoted to columns?
   - Example: `positions.detail.quantity` → `positions.quantity`
   - Trade-off: Flexibility vs Query Performance

3. **Materialized Views**: Should we add materialized views for common aggregations?
   - Example: `trader_performance_summary` (aggregated from `decision_cycles` + `trades`)

---

## 8. Next Steps 📝

### Immediate Actions
1. [ ] Verify table definitions against actual SQL schema files
2. [ ] Map each table to corresponding Go structs
3. [ ] Validate JSONB contents match struct definitions
4. [ ] Add missing indexes if verified as needed

### Follow-up Reviews
1. [ ] Review cross-references with `07a-domain-model.puml`
2. [ ] Validate foreign key relationships with actual constraints
3. [ ] Check if migration files exist for all tables

---

## 9. Sign-off ✍️

- **Status**: ⚠️ Approved with Comments
- **Comments**:
  - Diagram is well-structured and comprehensive
  - Minor improvements suggested for indexes and retention policy
  - Code verification pending
- **Next Actions**:
  1. Code alignment verification (Priority: High)
  2. Add retention policy notes (Priority: Medium)
  3. Consider materialized views (Priority: Low)

---

## Change Log

| Date | Reviewer | Changes |
|------|----------|---------|
| 2025-11-18 | Architecture Team | Initial review - comprehensive analysis |
