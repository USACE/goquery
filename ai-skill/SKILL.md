---
name: ai-skill
description: specialized instructions for implementing AI-driven features and agentic workflows within the goquery ecosystem
category: ai-capabilities
source: usace
---

# AI-Agentic Skill: goquery AI Integration

## Overview

This skill provides specialized guidance for implementing AI-driven logic, agentic workflows, and automated data processing using the `goquery` framework. It leverages the structured patterns defined in the `goquery` core library to ensure AI-generated code is type-safe, secure, and idiomatic.

**Core Repository:** `github.com/usace/goquery`
**Reference Docs:** `ai-skill/references/documentation.md`

## When to Use This Skill

Use this workflow when the task involves:
- Integrating LLMs or AI agents into a Go application using `goquery`.
- Designing autonomous data retrieval pipelines (e.g., an agent deciding which `TableDataSet` to query).
- Automating database schema migrations or mapping via AI-generated struct definitions.
- Implementing "Self-Healing" data layers where AI interprets query errors and suggests fixes.

## Implementation Guidelines

### 1. Knowledge Retrieval (CRITICAL)
Before generating any code, you **MUST** read:
- `README.md` for general project structure.
- `ai-skill/references/documentation.md` for technical implementation details (Connection mapping, Struct tags, and the Enterprise Pattern).

### 2. Leveraging the Enterprise Pattern
When an AI agent needs to interact with the database, **DO NOT** suggest raw SQL strings. Instead, guide the implementation toward the `TableDataSet` pattern:
- Define the `TableDataSet` with a clear `StatementKey`.
- Use `.StatementKey()` to decouple the AI's intent from the specific SQL implementation.
- Ensure `TableFields` are used to drive auto-generated `INSERT` operations.

### 3. Security & Safety Guardrails (Strict Adherence)
When generating code for AI agents that interact with databases, you must enforce these rules:
- **SQL Injection:** Always use `.Params()` for user/agent-provided values. Never use `.Apply()` for values; only for identifiers.
- **Credential Safety:** Never suggest hardcoding credentials. Always use `RdbmsConfigFromEnv()` or secret management.
- **Error Handling:** When an agent encounters a database error, suggest using `.LogSql(true)` to debug the generated statement before attempting a retry.

### 4. Agentic Loop Example
When designing an agent that "searches" the database:
1. **Plan:** Agent identifies the required entity (e.g., `Product`).
2. **Retrieve:** Agent looks up the corresponding `TableDataSet` in the codebase.
3. **Execute:** Agent uses `.Select().DataSet(&ds).StatementKey("key").Fetch()`.
4. **Verify:** Agent checks `err` and uses the "Troubleshooting" section of the documentation to resolve issues.