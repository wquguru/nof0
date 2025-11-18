# PlantUML Architecture Review Guide

> **Quick Start**: How to review complex system diagrams systematically

---

## 🎯 Purpose

This directory contains the **comprehensive review framework** for evaluating PlantUML architecture diagrams in a complex system. Use these templates and processes to:

1. Ensure **consistency** across diagrams
2. Validate **architectural principles** (DDD, Clean Architecture, SOLID)
3. Catch **design flaws** early
4. Maintain **documentation quality**

---

## 📂 Directory Structure

```
review/
├── README.md                                  # This file - Quick start guide
├── TEMPLATE-single-diagram-review.md          # Template for reviewing individual diagrams
├── 07b-persistence-model-review.md            # Example: Sample individual review
├── cross-diagram-consistency.md               # Cross-diagram consistency checks
├── architecture-coherence-validation.md       # Architecture principle validation
└── EXECUTIVE-SUMMARY.md                       # Final review report
```

---

## 🚀 Quick Start (5-Step Process)

### Step 1: Review Individual Diagrams

Use the **template** for each diagram:

```bash
# 1. Copy template
cp TEMPLATE-single-diagram-review.md reviews/[diagram-name]-review.md

# 2. Fill in metadata
# Edit: Diagram name, reviewer, date

# 3. Complete checklist
# Work through each section systematically

# 4. Verify against code
grep -r "EntityName" ../../../pkg/
```

**See Example**: `07b-persistence-model-review.md`

---

### Step 2: Check Cross-Diagram Consistency

Open `cross-diagram-consistency.md` and:

1. **Fill Entity Matrix** - Track entity names across diagrams
2. **Verify Relationships** - Check if relationships match
3. **Trace Data Flows** - Follow data from source to sink
4. **Validate States** - Ensure state definitions are consistent

**Example Check**:
```markdown
| Entity | Domain Model | Persistence | Lifecycle | Sequence |
|--------|-------------|------------|-----------|----------|
| Trader | ✅ Aggregate | ❌ (runtime) | ✅ States | ✅ Used |
```

---

### Step 3: Validate Architecture Principles

Use `architecture-coherence-validation.md`:

1. **DDD Checks**:
   - Are aggregate roots clearly marked?
   - Are boundaries respected?
   - Entity vs Value Object distinction clear?

2. **Clean Architecture Checks**:
   - Do dependencies point inward?
   - Are ports/adapters separated?
   - Is domain independent of frameworks?

3. **SOLID Checks**:
   - Single Responsibility?
   - Interface Segregation?
   - Dependency Inversion?

**Quick Validation**:
```plantuml
' ✅ GOOD: Domain → Port → Adapter → External
Domain --> IPort : uses
Adapter ..|> IPort : implements
Adapter --> External : calls

' ❌ BAD: Domain → Adapter (coupling)
Domain --> Adapter : VIOLATION
```

---

### Step 4: Generate Summary Report

Fill in `EXECUTIVE-SUMMARY.md`:

1. **Overall Score** (out of 10)
2. **Key Findings**:
   - Critical issues (must fix)
   - Medium issues (should fix)
   - Minor issues (nice to have)
3. **Recommendations**
4. **Next Steps**

---

### Step 5: Track Progress

Use a tracking sheet:

| Diagram | Individual Review | Cross-Check | Sign-off |
|---------|------------------|------------|----------|
| 01-system-architecture | ⬜ | ⬜ | ⬜ |
| 02-trading-decision | ⬜ | ⬜ | ⬜ |
| 02a-config-core | ⬜ | ⬜ | ⬜ |
| ... | ... | ... | ... |

---

## 🔍 Review Checklists

### Individual Diagram Review

**Metadata** (1 min):
- [ ] Title clear?
- [ ] Theme consistent?
- [ ] Legend present?

**Content** (10 min):
- [ ] All entities present?
- [ ] Relationships defined?
- [ ] Colors meaningful?
- [ ] Text readable?

**Accuracy** (15 min):
- [ ] Names match code?
- [ ] Types correct?
- [ ] Relationships match implementation?

**Domain-Specific** (10 min):
- [ ] For Config: Defaults shown?
- [ ] For Domain: Aggregates marked?
- [ ] For Persistence: PKs/FKs marked?
- [ ] For Lifecycle: States complete?

**Total**: ~35 min per diagram

---

### Cross-Diagram Consistency

**Entity Names** (15 min):
- [ ] Build entity matrix
- [ ] Check for aliases/typos
- [ ] Verify consistency

**Relationships** (20 min):
- [ ] Map relationships across diagrams
- [ ] Check for contradictions
- [ ] Verify cardinality

**Data Flow** (30 min):
- [ ] Trace end-to-end flows
- [ ] Check persistence points
- [ ] Verify transformations

**States** (15 min):
- [ ] List all states per entity
- [ ] Check for missing transitions
- [ ] Verify state names

**Total**: ~80 min (one-time)

---

### Architecture Validation

**DDD** (30 min):
- [ ] Aggregate roots identified?
- [ ] Boundaries clear?
- [ ] Invariants protected?

**Clean Architecture** (20 min):
- [ ] Dependency direction correct?
- [ ] Ports defined?
- [ ] Adapters implement ports?

**SOLID** (40 min):
- [ ] SRP: Single responsibility?
- [ ] OCP: Open for extension?
- [ ] LSP: Substitution works?
- [ ] ISP: No fat interfaces?
- [ ] DIP: Depend on abstractions?

**Total**: ~90 min (one-time)

---

## 🛠️ Tools & Scripts

### Diagram Rendering

```bash
# Render all diagrams
plantuml docs/flows/*.puml -o output/

# Render with specific format
plantuml diagram.puml -tpng
plantuml diagram.puml -tsvg
```

### Entity Extraction

```bash
# List all classes/entities
grep -oh "class [A-Za-z]*" docs/flows/*.puml | sort | uniq

# List all relationships
grep -E "(-->|--|\.\.>)" docs/flows/*.puml
```

### Consistency Checking

```bash
# Find entity name variations
grep -r "Trader" docs/flows/*.puml | grep -v "TraderConfig"

# Check for TODO/FIXME markers
grep -r "TODO\|FIXME\|XXX" docs/flows/*.puml
```

### Code Alignment

```bash
# Find struct definition
grep -r "type Trader struct" pkg/

# Find interface definition
grep -r "type IExchangeProvider interface" pkg/

# Find YAML config
grep -r "decision_interval" etc/
```

---

## 📊 Review Metrics

Track these metrics over time:

| Metric | Target | Current |
|--------|--------|---------|
| Diagrams reviewed | 100% | [Fill] |
| Consistency score | >90% | [Fill] |
| Architecture score | >8/10 | [Fill] |
| Critical issues | 0 | [Fill] |
| Time to review | <8 hours | [Fill] |

---

## 🎓 Best Practices

### DO ✅
- Review diagrams **bottom-up** (persistence → domain → flows → architecture)
- Use the **template** for consistency
- **Cross-check** entities across diagrams
- **Verify** against actual code
- Document **findings** with evidence
- Prioritize issues by **impact**

### DON'T ❌
- Don't skip individual reviews to save time
- Don't assume diagrams match code without verification
- Don't ignore minor inconsistencies (they accumulate)
- Don't review in isolation (pair review is better)
- Don't forget to re-review after fixes

---

## 🔄 Review Frequency

| Trigger | Review Type | Effort |
|---------|------------|--------|
| **New diagram added** | Individual review only | 35 min |
| **Major refactoring** | Full cross-diagram review | 3 hours |
| **Quarterly** | Full architecture review | 8 hours |
| **Before release** | Complete validation | 8 hours |

---

## 🚨 Common Issues & Solutions

### Issue: Entity Name Inconsistency

**Symptom**: `Trader` in one diagram, `TraderEntity` in another

**Solution**:
1. Build entity matrix (Step 2)
2. Choose canonical name
3. Update all diagrams
4. Commit with clear message

### Issue: Relationship Contradiction

**Symptom**: Diagram A shows 1:N, Diagram B shows 1:1

**Solution**:
1. Check actual code implementation
2. Identify correct cardinality
3. Update incorrect diagram(s)
4. Add note explaining the relationship

### Issue: Missing Error Paths

**Symptom**: Sequence diagrams only show happy path

**Solution**:
```plantuml
alt Success
    A -> B: Request
    B --> A: Response
else Error
    A -> B: Request
    B --> A: Error
    A -> A: Log error
end
```

---

## 📚 Reference Materials

- **DDD Blue Book**: Eric Evans, Domain-Driven Design
- **Clean Architecture**: Robert C. Martin
- **PlantUML Guide**: https://plantuml.com/guide
- **Project CLAUDE.md**: `docs/flows/CLAUDE.md`

---

## 🙋 FAQ

**Q: How long should a full review take?**
A: For 11 diagrams:
- Individual reviews: 11 × 35min = 6.5 hours
- Cross-checks: 1.5 hours
- Architecture validation: 1.5 hours
- Report: 0.5 hours
- **Total**: ~10 hours

**Q: Should we review every change?**
A: No, only:
- New diagrams (individual review)
- Major refactors (full review)
- Quarterly (full review)

**Q: Who should review?**
A: Ideally:
- **Primary**: Architect + Senior Engineer (pair review)
- **Secondary**: Domain Expert (business logic check)
- **Final**: Team Lead (sign-off)

**Q: How to handle disagreements?**
A:
1. Document both viewpoints
2. Check code implementation
3. Consult CLAUDE.md principles
4. Team decision (majority vote)
5. Document decision rationale

---

## 📝 Change Log

| Date | Reviewer | Changes |
|------|----------|---------|
| 2025-11-18 | Architecture Team | Initial review framework created |
| 2025-11-18 | Architecture Team | Completed sample review (07b) |
| 2025-11-18 | Architecture Team | Generated executive summary |

---

## ✍️ Contacts

- **Questions**: Ask in #architecture Slack channel
- **Issues**: File issue in GitHub repo
- **Updates**: Submit PR with updated review files

---

**Happy Reviewing! 🎉**

_Remember: Good diagrams prevent bad code. Time spent reviewing is time saved debugging._
