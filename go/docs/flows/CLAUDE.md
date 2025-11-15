# ARCHITECTURAL GUIDANCE FOR LLM-DRIVEN TRADING SYSTEM

## PURPOSE

This document provides **core architectural principles** for implementing a Go-based LLM-driven trading system. All detailed specifications exist in PlantUML diagrams (`.puml` files) in this directory - **diagrams are the authoritative source**.

---

## DIAGRAM INDEX

```
01-system-architecture.puml      - Component boundaries
02-trading-decision-cycle.puml   - Full decision flow
03-executor-decision-flow.puml   - LLM prompt → validation
04-order-execution-flow.puml     - Decision → order placement
05-risk-management-guards.puml   - Risk gate enforcement
06-data-ingestion-flow.puml      - Market data pipeline
07-trader-lifecycle.puml         - Trader state machine
08-component-dependencies.puml   - Dependency graph
09a-ideal-domain-model.puml      - Target domain model (DDD)
09b-ideal-persistence-model.puml - Target database schema
10-ideal-entity-lifecycle.puml   - Entity state management
```

---

## CORE PRINCIPLES

### P1: Diagram Authority
**Rule**: All implementation MUST match PlantUML diagrams. Component boundaries, dependencies, and flows defined in diagrams 01-10 are non-negotiable.

**Go-Zero Alignment**:
- PUML `package` → Go package with interface
- `<<core>>` → Business logic in `pkg/`
- `<<provider>>` → External adapters in `pkg/*/provider.go`
- `<<data>>` → Repository pattern in `pkg/repo/`

**Reference**: `@startuml 01-system-architecture.puml`, `@startuml 08-component-dependencies.puml`

---

### P2: Dependency Injection via ServiceContext
**Pattern**: All dependencies declared as interfaces, injected through `internal/svc/ServiceContext`.

**Rules**:
1. NO global variables for stateful dependencies
2. ALL external dependencies MUST have interface abstractions
3. Constructor functions receive all dependencies as parameters
4. Test doubles injected via same interfaces

**Reference**: `@startuml 08-component-dependencies.puml`

---

### P3: Domain-Driven Design (DDD)
**Reference**: `@startuml 09a-ideal-domain-model.puml`

**Key Concepts**:
- **Aggregate Root**: `Trader` - contains `TraderRuntimeState`, `PerformanceMetrics`, `CashAllocation`
- **Entity**: Has identity (ID field) - e.g., `Trader`, `Position`, `DecisionCycle`
- **Value Object**: No identity, immutable - e.g., `RiskParameters`, `ExecGuards`, `Symbol`
- **Aggregate Boundary**: External access ONLY through aggregate root methods, NEVER direct state mutation

**Immutability Contract**:
```go
// ✓ Value Object (immutable)
type RiskParameters struct { MaxPositions int }
func (r RiskParameters) WithMaxPositions(n int) RiskParameters { ... }

// ✗ FORBIDDEN: Mutable value object
func (r *RiskParameters) SetMaxPositions(n int) { r.MaxPositions = n }
```

---

### P4: Persistence Model
**Reference**: `@startuml 09b-ideal-persistence-model.puml`

**Table Types** (color-coded in PUML):
- **Configuration** (#E3F2FD): Immutable configs - `trader_config`, `models`, `symbols`
- **State** (#FFF3E0): Current state - `accounts`
- **Event Log** (#E8F5E9): Append-only - `positions`, `trades`, `decision_cycles`, `conversation_messages`
- **Snapshot** (#FCE4EC): Time-series - `account_snapshots`, `market_metrics`

**JSONB Pattern**:
- Store domain-specific details in `detail` JSONB column
- Enables flexible schema without excessive JOINs
- Index provider/type columns for queries

**Event Sourcing**:
- Every decision cycle → new `decision_cycles` row
- Enables replay, audit trail, state reconstruction

---

### P5: State Machine Rigor
**Reference**: `@startuml 07-trader-lifecycle.puml`, `@startuml 10-ideal-entity-lifecycle.puml`

**Enforcement**:
1. ALL state changes MUST go through transition guard methods
2. NO direct state field assignment
3. State changes MUST be logged to audit log

**Allowed Transitions**:
```
Running → {Paused, Stopped}
Paused  → {Running, Stopped}
Stopped → {} (terminal)
```

---

### P6: Risk Guard Layering
**Reference**: `@startuml 05-risk-management-guards.puml`

**Sequential Execution Order** (fail-fast):
1. Max positions limit
2. Margin usage check
3. Liquidity threshold
4. Position value bands (BTC/ETH vs Alt)
5. Symbol cooldown check
6. Performance gating (Sharpe-based pause)

**Feature Toggles**: Each guard has `Enable*Guard` flag (nil = default true)

---

### P7: Decision Cycle Orchestration
**Reference**: `@startuml 02-trading-decision-cycle.puml`

**Flow**:
1. **Scheduling**: Ticker every N seconds
2. **Context Building**: Aggregate account, positions, market snapshots, performance
3. **LLM Decision**: Render prompt → call LLM → validate schema → validate business rules
4. **Execution**: Sort (close first) → cap new opens → execute with guards
5. **Persistence**: Update runtime state → record decision cycle

---

### P8: Executor Pattern (LLM Integration)
**Reference**: `@startuml 03-executor-decision-flow.puml`

**Pipeline**:
1. **Prompt Rendering**: Jet template with context variables
2. **LLM Call**: OpenAI/Anthropic API with structured output
3. **Schema Validation**: JSON Schema validation (`gojsonschema`)
4. **Business Validation**: Risk-reward ratio, confidence threshold, leverage limits

---

### P9: Order Execution Flow
**Reference**: `@startuml 04-order-execution-flow.puml`

**Processing Order**:
1. Sort decisions: **Close positions FIRST**
2. Cap new positions per cycle
3. For close: Exchange close → record cooldown → log position event
4. For open: Enforce risk guards → update leverage → get price → place order → log event

---

### P10: Data Ingestion Pattern
**Reference**: `@startuml 06-data-ingestion-flow.puml`

**Market Data Caching**:
- Use go-zero `cache.Cache` for snapshots
- TTL-based expiration (default 5s)
- Cache key: `snapshot:<symbol>`

**Symbol Universe**:
- Periodic refresh from exchange
- Filter: perpetuals, non-delisted
- Persist to `symbols` table

---

## GO-ZERO PATTERNS

**Configuration**: YAML-based with `config.Config`
**Logging**: Structured logging with `logx.WithContext(ctx).Infow(...)`
**Error Handling**: Always wrap with context - `fmt.Errorf("...: %w", err)`
**Graceful Shutdown**: Signal handlers for `SIGINT`/`SIGTERM`

---

## TESTING STRATEGIES

**Unit Tests**: Mocked dependencies via interfaces
**Integration Tests**: Real database + mocked external providers
**Table-Driven Tests**: For validators and business rules

---

## ANTI-PATTERNS TO AVOID

1. **Bypassing Aggregate Boundaries**: Direct state mutation
2. **Mutable Value Objects**: Pointer receivers on value objects
3. **God Services**: Manager doing everything
4. **Global State**: Global variables for stateful dependencies
5. **Missing Error Context**: Silent error swallowing
6. **Tight Coupling**: Direct HTTP calls in business logic

---

## IMPLEMENTATION CHECKLIST

**Before Implementation**:
- [ ] Identify which PUML diagram(s) describe the feature
- [ ] Verify component matches diagram boundaries
- [ ] Check dependencies flow according to PUML arrows

**During Implementation**:
- [ ] Entity/Value Object distinction clear
- [ ] Aggregate boundary respected
- [ ] State transitions via state machine methods
- [ ] Table type matches PUML schema
- [ ] Guard execution order matches PUML

**After Implementation**:
- [ ] Unit tests with mocked dependencies
- [ ] Integration tests with real database
- [ ] Error cases covered
- [ ] Structured logging added

---

## MAINTENANCE PROTOCOL

**Modifying Diagrams**:
1. Update PUML diagram(s) first
2. Document changes in git commit
3. Update implementation to match
4. Update tests

**Adding Features**:
1. Identify affected PUML diagram(s)
2. Sketch new components/flows in PUML
3. Review before implementation
4. Implement according to diagram

**Debugging**:
1. Trace flow in sequence diagrams (02-06)
2. Verify state transitions (07, 10)
3. Check dependency graph (08)
4. Confirm persistence matches schema (09b)

---

## QUICK REFERENCE

| Concern | Diagram |
|---------|---------|
| Domain Model | `09a-ideal-domain-model.puml` |
| Database Schema | `09b-ideal-persistence-model.puml` |
| Decision Flow | `02-trading-decision-cycle.puml` |
| Risk Guards | `05-risk-management-guards.puml` |
| State Machine | `07-trader-lifecycle.puml`, `10-ideal-entity-lifecycle.puml` |
| Dependencies | `08-component-dependencies.puml` |

---

**Diagrams are authoritative. When in doubt, consult PUML.**
