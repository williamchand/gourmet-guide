# Claude Collaboration Prompt (GourmetGuide)

Use this operating mode when reviewing plans or code for this repository.

## Review Meta-Instruction
Review thoroughly before making code changes.
For every issue or recommendation:
- Explain concrete tradeoffs.
- Give an opinionated recommendation.
- Ask for user input before choosing a direction.

## Engineering Preferences
- DRY is important — flag repetition aggressively.
- Testing is non-negotiable — prefer more tests over fewer.
- Aim for "engineered enough" solutions.
- Bias toward handling edge cases thoughtfully.
- Prefer explicit over clever.

## Review Framework
### 1) Architecture Review
Evaluate:
- Overall system design and component boundaries
- Dependency graph and coupling concerns
- Data flow and bottlenecks
- Scalability and single points of failure
- Security architecture (auth, data access, API boundaries)

### 2) Code Quality Review
Evaluate:
- Module structure and organization
- DRY violations (aggressively)
- Error handling and edge-case gaps
- Technical debt hotspots
- Over- vs under-engineering relative to preferences

### 3) Test Review
Evaluate:
- Coverage gaps (unit/integration/e2e)
- Assertion quality
- Missing edge-case coverage
- Untested failure modes

### 4) Performance Review
Evaluate:
- N+1 query risks
- Memory concerns
- Caching opportunities
- High-complexity paths

## Required Interaction Pattern
Before starting, ask user to choose one:
1. **BIG CHANGE**: Interactive, one section at a time with at most 4 top issues per section.
2. **SMALL CHANGE**: Interactive, one question per review section.

For each stage:
- Number issues (`Issue 1`, `Issue 2`, ...).
- Provide options labeled with letters (`A`, `B`, `C`), including "do nothing" when reasonable.
- Make the recommended option appear first.
- Include pros/cons for each option.
- Explicitly ask user which numbered issue + option letter they choose.
