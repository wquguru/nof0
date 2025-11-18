# Single Diagram Review Template

> **Diagram:** [diagram-name].puml
> **Reviewer:** [Your Name]
> **Date:** [YYYY-MM-DD]
> **Status:** 🔴 Not Started | 🟡 In Progress | 🟢 Completed

---

## 1. Metadata Check ✅

- [ ] **Title**: Clear and descriptive?
- [ ] **Theme**: Consistent with other diagrams?
- [ ] **Legend**: Present and helpful?
- [ ] **Notes**: Key decisions explained?

---

## 2. Content Review 📋

### 2.1 Completeness
- [ ] All required entities present?
- [ ] All relationships defined?
- [ ] Missing elements identified?

**Missing Elements:**
```
- [List any missing elements]
```

### 2.2 Accuracy
- [ ] Entity names match codebase?
- [ ] Relationships match implementation?
- [ ] Data types correct?

**Discrepancies Found:**
```
- [List any discrepancies with code]
```

### 2.3 Clarity
- [ ] Easy to understand at first glance?
- [ ] Logical grouping of components?
- [ ] Colors used effectively?
- [ ] Text readable (no overlaps)?

**Clarity Issues:**
```
- [List any clarity issues]
```

---

## 3. Domain-Specific Checks 🎯

### For Configuration Diagrams (02a-config-*.puml)
- [ ] All config fields documented?
- [ ] Default values shown?
- [ ] Environment variables marked?
- [ ] Validation rules noted?

### For Domain Model Diagrams (07a-domain-model.puml)
- [ ] Aggregate roots clearly marked?
- [ ] Entities vs Value Objects distinguished?
- [ ] Aggregate boundaries defined?
- [ ] Invariants documented?

### For Persistence Diagrams (07b-persistence-model.puml)
- [ ] Primary keys marked?
- [ ] Foreign keys marked?
- [ ] Indexes documented?
- [ ] JSONB fields explained?
- [ ] Table types color-coded (Config/State/Event/Snapshot)?

### For State Machine Diagrams (08-entity-lifecycle.puml)
- [ ] All states defined?
- [ ] Transitions complete?
- [ ] Entry/Exit actions documented?
- [ ] Guard conditions shown?
- [ ] Terminal states marked?

### For Sequence Diagrams (02-06-*.puml)
- [ ] Actors/participants clear?
- [ ] Message flow logical?
- [ ] Alternative paths shown?
- [ ] Error handling included?
- [ ] Performance considerations noted?

---

## 4. Code Alignment Check 🔗

**Files to verify against:**
```
- pkg/[module]/[file].go
- internal/svc/servicecontext.go
- etc/[config].yaml
```

**Verification Results:**

| Diagram Element | Code Location | Status | Notes |
|----------------|---------------|--------|-------|
| Entity X       | pkg/foo/x.go  | ✅ Match | -     |
| Field Y        | pkg/foo/x.go  | ⚠️ Mismatch | Type different |
| Relationship Z | -             | ❌ Missing | Not implemented |

---

## 5. Issues Found 🐛

### Critical Issues (Must Fix)
1. [Issue description]
   - **Location**: [Where in diagram]
   - **Impact**: [Why critical]
   - **Suggested Fix**: [How to fix]

### Medium Issues (Should Fix)
1. [Issue description]

### Minor Issues (Nice to Have)
1. [Issue description]

---

## 6. Recommendations 💡

### Improvements
- [ ] Recommendation 1
- [ ] Recommendation 2

### Questions for Team
1. Question 1?
2. Question 2?

---

## 7. Sign-off ✍️

- **Status**: ✅ Approved | ⚠️ Approved with Comments | ❌ Rejected
- **Comments**: [Final comments]
- **Next Actions**: [What needs to be done]

---

## Change Log

| Date | Reviewer | Changes |
|------|----------|---------|
| YYYY-MM-DD | Name | Initial review |
