---
title: "Workflow Failure: {{ env.workflow_name }}"
labels: bug, ci-failure
assignees: ""
---

## Workflow Failure Alert

A critical workflow has failed and requires immediate attention.

### Details

- **Workflow**: {{ env.workflow_name }}
- **Run ID**: {{ env.run_id }}
- **Repository**: {{ env.repo_name }}
- **Triggered by**: {{ env.actor }}
- **Failure Time**: {{ date | date('YYYY-MM-DD HH:mm:ss') }}

### Next Steps

1. Check the [workflow logs]({{ env.workflow_url }}) for detailed error information
2. Investigate the root cause of the failure
3. Fix the issue and re-run the workflow
4. Close this issue once resolved

@Perun-Engineering/maintainers
