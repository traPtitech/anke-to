// This file exists only to satisfy organization-level CodeQL JavaScript
// scanning in a Go-only repository. Some org-level CodeQL configurations
// fail the analysis when they include JavaScript/TypeScript but the
// repository contains no JS/TS source files at all.

"use strict";

function codeqlPlaceholder() {
  return "codeql-placeholder";
}

module.exports = { codeqlPlaceholder };
