# AI-DLC and Spec-Driven Development

Kiro-style Spec Driven Development implementation on AI-DLC (AI Development Life Cycle)

## Project Context

### Paths
- Steering: `.kiro/steering/`
- Specs: `.kiro/specs/`

### Steering vs Specification

**Steering** (`.kiro/steering/`) - Guide AI with project-wide rules and context
**Specs** (`.kiro/specs/`) - Formalize development process for individual features

### Active Specifications
- Check `.kiro/specs/` for active specifications
- Use `/kiro:spec-status [feature-name]` to check progress

## Development Guidelines
- Think in English, generate responses in Japanese. All Markdown content written to project files (e.g., requirements.md, design.md, tasks.md, research.md, validation reports) MUST be written in the target language configured for this specification (see spec.json.language).

## Minimal Workflow
- Phase 0 (optional): `/kiro:steering`, `/kiro:steering-custom`
- Phase 1 (Specification):
  - `/kiro:spec-init "description"`
  - `/kiro:spec-requirements {feature}`
  - `/kiro:validate-gap {feature}` (optional: for existing codebase)
  - `/kiro:spec-design {feature} [-y]`
  - `/kiro:validate-design {feature}` (optional: design review)
  - `/kiro:spec-tasks {feature} [-y]`
- Phase 2 (Implementation): `/kiro:spec-impl {feature} [tasks]`
  - `/kiro:validate-impl {feature}` (optional: after implementation)
- Progress check: `/kiro:spec-status {feature}` (use anytime)

## Development Rules
- 3-phase approval workflow: Requirements → Design → Tasks → Implementation
- Human review required each phase; use `-y` only for intentional fast-track
- Keep steering current and verify alignment with `/kiro:spec-status`
- Follow the user's instructions precisely, and within that scope act autonomously: gather the necessary context and complete the requested work end-to-end in this run, asking questions only when essential information is missing or the instructions are critically ambiguous.

## Steering Configuration
- Load entire `.kiro/steering/` as project memory
- Default files: `product.md`, `tech.md`, `structure.md`
- Custom files are supported (managed via `/kiro:steering-custom`)

## Subagent-based Quality Assurance

This project uses Claude Code's subagent functionality for specialized validation tasks.

### Available Subagents (`.claude/agents/`)

1. **spec-validator** - Requirements & specification validation
   - Validates completeness and clarity of specifications
   - Ensures alignment with steering documents
   - Checks feasibility across all SDD phases
   - Usage: "Use spec-validator to validate the requirements for {feature}"

2. **design-reviewer** - Technical design review
   - Conducts architectural soundness reviews
   - Ensures requirements traceability
   - Identifies implementation risks
   - Provides GO/NO-GO decisions with rationale
   - Usage: "Use design-reviewer to review the design for {feature}"

3. **implementation-auditor** - Implementation verification
   - Validates implementation against requirements and design
   - Audits TDD compliance and test coverage
   - Checks code quality and security standards
   - Verifies task completion accuracy
   - Usage: "Use implementation-auditor to audit the implementation for {feature}"

### When to Use Subagents

- **After `/kiro:spec-requirements`**: Use spec-validator for requirements review
- **After `/kiro:spec-design`**: Use design-reviewer (alternative to `/kiro:validate-design`)
- **After `/kiro:spec-impl`**: Use implementation-auditor (alternative to `/kiro:validate-impl`)

Subagents provide isolated, focused validation with specialized expertise, keeping the main conversation uncluttered while ensuring comprehensive quality assurance.
