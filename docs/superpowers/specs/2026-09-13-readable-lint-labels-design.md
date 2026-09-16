---
status: active
last-verified: 2026-09-13
---

# Reader-friendly lint labels

The #870 lint report currently emits internal rule numbers and labels such as `R6`, `devlog-comment`, and `eta-expansion`. A reader should understand every finding, summary, skipped rule, and ignored-gate note without opening source code or knowing the rule catalog.

Keep EDN `:kind` values and `--gate` tokens stable. Text output uses plain-language display labels; it must not rely on `R1`–`R6` to identify a rule. A built-in map covers comment kinds. The data-driven code-rule catalog carries an optional `:label` field; existing/custom entries without one fall back to space-separated words from `:kind`, preserving catalog-only rule additions. All current catalog entries get explicit labels, including `eta-expansion` → “argument-forwarding wrapper” and `composable-accessor` → “nested sequence accessor”. Opaque built-in kinds become `restatement` → “comment repeats code”, `comment-divider` → “decorative section divider”, `comment-density-outlier` → “unusually comment-heavy definition”, and `devlog-comment` → “development-note phrase”. Evidence remains next to each finding and explains the match.

Tests cover every built-in kind and emitted branch: density skipped/calibrated; churn skipped/empty/finding/zero baseline; ignored gates; comment and code-verbosity findings. Existing semantic tests that search raw text labels move to display assertions or EDN-kind assertions. No detection threshold, gate eligibility, or exit status changes.
