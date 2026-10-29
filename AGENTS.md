# Agent Documentation - Project Setup & Architecture

This document provides a guide for AI agents working on `chrono` (the Go port of chrono-node) to quickly understand the project context, layout, testing commands, and implementation constraints.

## Project Overview
`chrono` is a natural language date parser. It extracts dates, times, and relative date expressions from free text. This repository is a Go module, `github.com/wanasit/chrono-go`, that reproduces chrono-node 2.10.1 exactly, including what the original took from the JavaScript runtime (UTF-16 strings, ECMAScript regular expressions, the `Date` object in the process time zone).

---

## Technical Stack & Setup
- **Language**: Go 1.27, standard library only (no module dependencies)
- **Testing Framework**: `go test`

### Running Tests
To build, vet and run the full test suite, execute:
```bash
go build ./... && go vet ./... && go test ./...
```

---

## Directory Structure & Architecture

```
chrono.go                 # Public API of the root package (types, Chrono, Casual/Strict, Parse/ParseDate)
locales/                  # Public locale packages (Casual, Strict, Parse, ParseDate, configuration constructors)
│   ├── de/
│   ├── en/
│   ├── zh/
│   │   ├── hans/
│   │   └── hant/
│   └── ...
internal/
├── core/                 # Chrono engine, results, time zones, calculations, common parsers and refiners
├── en/, de/, ...         # Language-specific configurations, parsers, and refiners
├── js/                   # JavaScript string, number and sorting semantics (UTF-16, case mapping, parseInt, ...)
├── jsregex/              # ECMAScript regular expression engine
├── jsdate/               # JavaScript Date semantics on Go's time zone data
└── uprops/               # Unicode property and case-folding tables used by jsregex
```

### Core Concepts
1. **Parsers (`Parser`)**:
   Return an ECMAScript regular expression from `Pattern()`, and process each match in `Extract(context, match)`, returning components or a `ParsingResult`.
2. **Refiners (`Refiner`)**:
   Post-process and merge/filter the parsed results (e.g., merging separate date and time components, or resolving overlapping date ranges).
3. **ParsingComponents (`ParsingComponents`)**:
   Stores date and time elements (e.g., year, month, day, hour). Elements are marked as either explicitly assigned (`Assign`) or implied (`Imply`).

---

## Key Context & Caveats

1. **Faithfulness**:
   Results must match chrono-node byte for byte. Keep each parser and refiner as its own value, in the original's configuration order, and keep JavaScript semantics (UTF-16 indices, ECMAScript regex flags, `Date` local time, object key order) by using the helpers in `internal/js`, `internal/jsregex` and `internal/jsdate`.

2. **Locale Symmetry & Chinese (zh) Special Case**:
   - Simplified (`zh/hans`) and Traditional (`zh/hant`) components stay exactly synchronized; implement a change in both.
   - English (`en`) is the most advanced locale. Other locales are modeled after English but often lag behind.

3. **Relative Date Arithmetic**:
   Use the duration helpers in `internal/core` (`AddDuration`, `ReverseDuration`). When handling "ago/before" expressions, calculate the positive duration first, then reverse it.

4. **Testing Conventions**:
   - Always pass a reference instant so relative expressions are deterministic.
   - Corpus tests replay recorded outputs of chrono-node from `testdata/`.

---

## Git & Commit Message Format

Follow the Conventional Commits specification:
```
<type>(<scope>): <description>
```
- Example: `feat(zh): support past-tense relative time expressions (e.g. "1小时前")`
- Example: `fix(en): support 'of' connector in weekday postfix`

Run `gofmt -w .` before committing.

---

## AI Agent Tool & Command Guidelines

To ensure secure, auditable, and standard operations:
- **Network Requests**: Do not run arbitrary commands to execute network calls.
- **Verification**: Always run `go test ./...` inside the workspace directory to verify behaviors.
