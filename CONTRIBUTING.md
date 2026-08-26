<!--
SPDX-FileCopyrightText: 2026 Deutsche Telekom AG

SPDX-License-Identifier: CC0-1.0
-->

# Contributing

All contributors must follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Build and test

Install these tools to run the full local check suite:

- The Go version declared in [`go.mod`](go.mod)
- [golangci-lint](https://golangci-lint.run/docs/welcome/install/) 2.13.1 or newer
- [REUSE](https://reuse.software/dev/#install) 6.0.0 or newer
- [govulncheck](https://go.dev/doc/tutorial/govulncheck), installed with
  `go install golang.org/x/vuln/cmd/govulncheck@latest`

Build and test the service:

```bash
make build
make test-unit
```

Run all source-level CI checks in sequence:

```bash
make check
```

`make check` runs formatting, lint, REUSE, build, race-enabled unit tests, and
govulncheck. Govulncheck may access the Go vulnerability database. The build and
test targets create ignored local artifacts. In its default source scan,
govulncheck exits unsuccessfully when the code calls a known vulnerable symbol,
including when no fix is available. It reports imported vulnerabilities that
are not called without failing the scan.

## Git hooks

The repository includes an optional [Lefthook](https://lefthook.dev/)
configuration. The hooks use `gofmt` from the Go installation and REUSE from
the build and test requirements above. Install these additional tools before
enabling the hooks:

- [Lefthook](https://lefthook.dev/#how-to-install-lefthook) 2.0.4 or newer
- [Gitleaks](https://github.com/gitleaks/gitleaks#installing) 8.20.0 or newer
- [committed](https://github.com/crate-ci/committed#install)

The hooks use tools from `PATH`. They do not install or upgrade tools. CI uses
the tools supplied by its GitHub Actions and remains the authoritative
validation gate.

Enable the hooks once per clone:

```bash
make hooks
```

The setup target checks that each required command exists, then runs
`lefthook install`.

| Hook | Checks |
| --- | --- |
| `pre-commit` | Check Go formatting when staged Go files change. Check whole-repository REUSE compliance and scan staged content with Gitleaks on every commit. The jobs run in parallel. |
| `commit-msg` | Validate Conventional Commits with `committed`, except during merges and rebases. |

The checks do not rewrite tracked source files. Run `make fmt` to fix Go
formatting. Gitleaks scans staged content only and has no matching CI job.
Build, test, lint, and govulncheck remain explicit `make check` and CI tasks
rather than Git hook tasks.

To bypass installed hooks for one command, set `LEFTHOOK=0` or use Git's
`--no-verify` option:

```bash
LEFTHOOK=0 git commit ...
git commit --no-verify ...
```

Use a bypass only when necessary. CI still runs its required checks.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/). The accepted
types match [`release.config.js`](release.config.js):

`feat`, `fix`, `build`, `chore`, `ci`, `docs`, `perf`, `refactor`, `revert`,
`style`, and `test`.

Scopes are optional. `committed.toml` does not impose subject or body length
limits beyond the Conventional Commits structure.
