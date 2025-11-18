# Cross-Diagram Consistency Check

> **Purpose**: Ensure all diagrams tell a coherent story
> **Date**: 2025-11-18
> **Status**: 🟡 In Progress

---

## 1. Entity Name Consistency Matrix 📊

### 1.1 Core Entities

| Entity Name | Config Diagrams | Domain Model | Persistence | Lifecycle | Sequence Diagrams | Status |
|-------------|----------------|--------------|-------------|-----------|-------------------|--------|
| **Trader** | ✅ TraderConfig | ✅ Aggregate | ❌ (runtime) | ✅ States | ✅ Participant | ✅ Consistent |
| **TraderConfig** | ✅ Entity | ✅ Entity | ✅ trader_config | - | - | ✅ Consistent |
| **Position** | - | ✅ Entity | ✅ positions | ✅ States | ✅ Participant | ✅ Consistent |
| **Decision** | - | ✅ Value Object | ❌ (ephemeral) | - | ✅ Data | ✅ Consistent |
| **DecisionCycle** | - | ✅ Entity | ✅ decision_cycles | ✅ States | ✅ Event | ✅ Consistent |
| **Executor** | ✅ ExecutorConfig | ✅ Service | - | - | ✅ Participant | ✅ Consistent |
| **Symbol** | - | ✅ Value Object | ✅ symbols | - | ✅ Data | ⚠️ VO vs Table |
| **RiskParameters** | ✅ Value Object | ✅ Value Object | ✅ JSONB field | - | - | ✅ Consistent |
| **ExecGuards** | ✅ Value Object | ✅ Value Object | ✅ JSONB field | - | - | ✅ Consistent |

**Legend:**
- ✅ Consistent: Entity name and semantics match across diagrams
- ⚠️ Needs Clarification: Minor inconsistency requiring documentation
- ❌ Inconsistent: Conflicting definitions

### 1.2 Findings

**Issue 1: Symbol - Value Object vs Table**
- **In Domain (07a)**: Value Object (immutable, no identity)
- **In Persistence (07b)**: Table with PK (has identity)
- **Resolution**: Acceptable - Reference data pattern
- **Action**: Add clarifying note in 07a-domain-model.puml

---

## 2. Relationship Consistency Check 🔗

### 2.1 Trader → Position Relationship

| Diagram | Relationship Type | Cardinality | Notes |
|---------|------------------|-------------|-------|
| 07a-domain-model | Aggregate manages | 1:N | ✅ |
| 07b-persistence | FK (trader_id) | 1:N | ✅ |
| 08-entity-lifecycle | State machine owns | - | ✅ |
| 02-trading-decision | Runtime tracking | - | ✅ |

**Status**: ✅ Consistent across all diagrams

### 2.2 Trader → TraderConfig Relationship

| Diagram | Relationship Type | Cardinality | Notes |
|---------|------------------|-------------|-------|
| 07a-domain-model | References | N:1 | ✅ |
| 07b-persistence | Not FK (config loaded at startup) | - | ✅ |
| 02-trading-decision | Loads from DB | - | ✅ |
| 01-system-architecture | Hydrates | - | ✅ |

**Status**: ✅ Consistent - Config is loaded once per trader lifecycle

### 2.3 Executor → Decision Relationship

| Diagram | Relationship Type | Cardinality | Notes |
|---------|------------------|-------------|-------|
| 07a-domain-model | Produces | 1:N | ✅ |
| 03-executor-decision | Output | - | ✅ |
| 02-trading-decision | GetFullDecision() | - | ✅ |

**Status**: ✅ Consistent

---

## 3. Data Flow Consistency 🌊

### 3.1 Decision Flow Trace

**Start**: User request → Manager
**Flow Diagram**: 02-trading-decision.puml

| Step | Diagram | Entity | Action | Next |
|------|---------|--------|--------|------|
| 1 | 02-trading-decision | Manager | buildContext() | → |
| 2 | 02-trading-decision | Executor | GetDecision(context) | → |
| 3 | 03-executor-decision | Executor | renderPrompt() | → |
| 4 | 03-executor-decision | LLMClient | Call(prompt) | → |
| 5 | 03-executor-decision | Executor | validateSchema() | → |
| 6 | 03-executor-decision | Executor | validateRules() | → |
| 7 | 02-trading-decision | Manager | executeDecisions() | → |
| 8 | 04-order-execution | Exchange | PlaceOrder() | → |
| 9 | 07b-persistence | DB | RecordTrade() | ✅ |

**Cross-Diagram Verification:**
- [ ] Step 1-2: Does `buildContext()` create the `Context` DTO defined in 07a? 🔍
- [ ] Step 3-4: Does `renderPrompt()` match template structure in 03? 🔍
- [ ] Step 8: Does `PlaceOrder()` match order types in 04? 🔍

**Action Required**: Trace each step against code implementation

### 3.2 Data Persistence Flow

| Source Diagram | Event | Persistence Diagram | Table | Status |
|----------------|-------|---------------------|-------|--------|
| 02-trading-decision | Decision cycle complete | 07b-persistence | decision_cycles | ✅ |
| 04-order-execution | Order placed | 07b-persistence | trades | ✅ |
| 04-order-execution | Position opened | 07b-persistence | positions | ✅ |
| 02-trading-decision | Performance updated | 07b-persistence | trader_runtime_state | ⚠️ |

**Issue**: `trader_runtime_state` table not explicitly shown in 07b-persistence
- **Action**: Add `trader_runtime_state` table to persistence diagram

---

## 4. State Consistency Check 🔄

### 4.1 Trader State Definitions

| Diagram | State Name | Valid? | Notes |
|---------|-----------|--------|-------|
| 07a-domain-model | Running, Paused, Stopped | ✅ | Enum definition |
| 08-entity-lifecycle | Running, Paused, Stopped, SharpeGated | ⚠️ | Extra state |
| 02-trading-decision | Running check | - | Implicit |

**Issue**: `SharpeGated` state in lifecycle but not in domain model enum
- **Resolution**: `SharpeGated` is a sub-state of `Paused` (indicated by `PauseUntil`)
- **Action**: Add clarifying note in 08-entity-lifecycle

### 4.2 Position State Definitions

| Diagram | State Name | Valid? | Notes |
|---------|-----------|--------|-------|
| 07a-domain-model | Open, Closed | ✅ | Enum |
| 07b-persistence | open, closed (lowercase) | ⚠️ | String values |
| 08-entity-lifecycle | Open, Closed | ✅ | States |

**Issue**: Case mismatch (Enum vs DB)
- **Resolution**: Go enum should map to lowercase DB values
- **Action**: Verify `PositionStatus.String()` method uses lowercase

---

## 5. Configuration Consistency 🛠️

### 5.1 Config Field Mapping

**Trace: ExecGuards**

| Field | 02a-config-trading | 07a-domain-model | 07b-persistence | Code Location |
|-------|-------------------|------------------|-----------------|---------------|
| MaxNewPositionsPerCycle | ✅ | ✅ | ✅ JSONB | 🔍 Verify |
| LiquidityThresholdUSD | ✅ | ✅ | ✅ JSONB | 🔍 Verify |
| EnableLiquidityGuard | ✅ | ✅ | ✅ JSONB | 🔍 Verify |
| SharpePauseThreshold | ✅ | ✅ | ✅ JSONB | 🔍 Verify |

**Action**: For each config field, verify:
```bash
# 1. Check struct definition
grep -r "SharpePauseThreshold" pkg/

# 2. Check YAML example
grep -r "sharpe_pause_threshold" etc/

# 3. Check database schema
grep -r "sharpe_pause_threshold" schema/
```

### 5.2 Default Value Consistency

| Field | Config Diagram Default | Code Default | Match? |
|-------|----------------------|--------------|--------|
| DecisionInterval | 3m | 🔍 Check | - |
| MaxPositions | 4 | 🔍 Check | - |
| MinConfidence | 75 | 🔍 Check | - |

---

## 6. Interface Consistency 🔌

### 6.1 Port Definitions

**From 01-system-architecture.puml:**

| Port (Interface) | Methods Documented? | Adapter Implements? | Used In Sequence? |
|------------------|--------------------|--------------------|-------------------|
| IExchangeProvider | ❌ | ✅ HyperliquidExchange | ✅ (02, 04) |
| IMarketProvider | ❌ | ✅ HyperliquidMarket | ✅ (02) |
| ILLMClient | ❌ | ✅ LLMClient | ✅ (03) |
| IRepository | ❌ | ✅ PostgresRepository | ✅ (02) |

**Issue**: Interface methods not documented in architecture diagram
- **Action**: Add method signatures to port definitions
- **Example**:
```plantuml
() "IExchangeProvider" as IExchange
note right of IExchange
  <b>Methods:</b>
  • GetAccount() Account
  • GetPositions() []Position
  • PlaceOrder(Order) Response
  • ClosePosition(Symbol) Response
end note
```

---

## 7. Consistency Issues Summary 🚨

### Critical (Must Fix)
1. ❌ **Missing `trader_runtime_state` table in 07b-persistence**
   - Impact: Incomplete persistence model
   - Fix: Add table definition with JSONB structure

### Medium (Should Fix)
1. ⚠️ **SharpeGated state not in domain model enum**
   - Impact: Confusion about valid states
   - Fix: Add note explaining it's derived from `PauseUntil`

2. ⚠️ **Interface methods not documented**
   - Impact: Unclear contracts
   - Fix: Add method signatures to 01-system-architecture

3. ⚠️ **Position state case mismatch**
   - Impact: Potential runtime bugs
   - Fix: Verify enum-to-string mapping

### Minor (Nice to Have)
1. ℹ️ **Symbol Value Object vs Table clarification**
   - Impact: Conceptual clarity
   - Fix: Add design note

---

## 8. Verification Checklist ✅

### Phase 1: Quick Wins
- [ ] Add `trader_runtime_state` to 07b-persistence.puml
- [ ] Add interface methods to 01-system-architecture.puml
- [ ] Add SharpeGated state note to 08-entity-lifecycle.puml
- [ ] Add Symbol design note to 07a-domain-model.puml

### Phase 2: Code Alignment
- [ ] Verify all config fields exist in structs
- [ ] Verify default values match between diagram/code
- [ ] Verify state enum values match DB constraints
- [ ] Verify interface implementations match ports

### Phase 3: End-to-End Tracing
- [ ] Trace decision flow from 02 → 03 → 04 → 07b
- [ ] Trace lifecycle from 08 → 02 → 04
- [ ] Trace data ingestion from 06 → 07b

---

## 9. Sign-off ✍️

- **Status**: 🟡 In Progress (60% complete)
- **Next Review**: After Phase 1 fixes applied
- **Reviewer Notes**: Overall consistency is good, minor gaps to address

