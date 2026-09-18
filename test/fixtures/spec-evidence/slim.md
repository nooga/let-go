<!-- test/fixtures/spec-evidence/slim.md -->
# Slim fixture

Slim assertion rows are evaluated cell by cell; each cell must be truthy. [R-slim]

<!-- evidence: @R-slim -->
| table | check | check |
|---|---|---|
|! catalog | (= 14 14) | (pos? 9) |
|! catalog | (string? "x") | (nil? nil) |
