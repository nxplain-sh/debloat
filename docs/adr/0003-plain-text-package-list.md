# 0003. Plain-text package list

Date: 2026-09-20

## Status

Accepted

## Context

The tool needs a curated list of package names to act on. Options considered:

- Embed the list in the binary as Go source.
- Store it as JSON/YAML/TOML with metadata (vendor, risk, description).
- Ship a plain text file, one package name per line, `#` for comments.

The list is highly device- and taste-dependent: what is safe to remove on one
ROM may break another, so the shipped list is a starting point users edit.

## Decision

Keep the default list in `packages.txt` as plain text: one package name per
line, `#` comments, blank lines ignored. Let `-f`/`--file` point at any other
file. A package is only acted on when its line is enabled, so opting out is
commenting it out.

## Consequences

- Lists are diffable, greppable, paste-able, and reviewable in a pull request;
  no schema or parser library is needed.
- Per-package metadata (risk notes, vendor) lives in comments and prose
  instead of structured fields; the tool cannot filter on it.
- The file format is user-owned: nothing stops a user from running the tool
  against a list that omits recovery-critical packages.
- Reserving the text format means future structured metadata must stay
  backwards compatible with it.
