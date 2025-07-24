---
title: "Security Alert: Vulnerabilities Detected"
labels: security, high-priority
assignees: ""
---

## Security Vulnerabilities Detected

The automated security scan has detected potential vulnerabilities in the project dependencies.

### Summary

- **Nancy Issues**: {{ env.NANCY_ISSUES }} potential vulnerabilities found
- **Go Vulnerability Check**: {{ env.VULN_ISSUES }} potential vulnerabilities found
- **Scan Date**: {{ date | date('YYYY-MM-DD HH:mm:ss') }}

### Next Steps

1. Review the detailed reports in the workflow artifacts
2. Prioritize vulnerabilities based on severity
3. Update affected dependencies or implement mitigations
4. Re-run the security scan to verify fixes

### Reports

- [Nancy Vulnerability Report]({{ env.GITHUB_SERVER_URL }}/{{ env.GITHUB_REPOSITORY }}/actions/runs/{{ env.GITHUB_RUN_ID }})
- [Go Vulnerability Check Report]({{ env.GITHUB_SERVER_URL }}/{{ env.GITHUB_REPOSITORY }}/actions/runs/{{ env.GITHUB_RUN_ID }})

@Perun-Engineering/maintainers
