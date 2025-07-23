# Branch Protection Rules

This document outlines the branch protection rules for the Terraform YAML Flattener Provider repository.

## Main Branch Protection

The `main` branch is protected with the following rules:

### Required Reviews
- At least 1 approval is required before merging
- Dismisses stale pull request approvals when new commits are pushed
- Requires review from Code Owners if CODEOWNERS file is present

### Required Status Checks
The following status checks must pass before merging:
- Build workflow (all platforms)
- Unit tests
- Integration tests
- Security scanning

### Branch Restrictions
- Force pushes are not allowed
- Branch deletion is not allowed
- Only administrators can bypass these restrictions in emergency situations

## Merge Settings

The repository uses the following merge settings:

- **Squash and merge** is the preferred merge method
- Pull request title and description are used for the commit message
- Commits are squashed into a single commit for a clean history
- Linear history is maintained

## Protected Files

The following files have additional protection:
- `.github/workflows/*` - Requires administrator approval
- `go.mod`, `go.sum` - Requires security review

## Automation

Branch protection rules are enforced through GitHub repository settings and cannot be modified without appropriate permissions.

## Exceptions

In rare cases, administrators may need to bypass these protections. All such actions should be documented and communicated to the team.
