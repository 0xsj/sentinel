"""Validate worker metadata. Does not execute tasks or validate source architecture."""
import json
from pathlib import Path
import sys

STATES = {"planned", "ready", "waiting", "in_progress", "implemented", "complete", "blocked"}

def validate(root, document):
    errors = []
    if document.get("schema_version") != 1:
        errors.append("schema_version must be 1")
    tasks = document.get("tasks")
    if not isinstance(tasks, list):
        return errors + ["tasks must be a list"]
    by_id = {}
    for task in tasks:
        if not isinstance(task, dict) or not isinstance(task.get("id"), str):
            errors.append("every task needs a string id")
            continue
        ident = task["id"]
        if ident in by_id:
            errors.append(f"duplicate task: {ident}")
        by_id[ident] = task
    def safe_path(value):
        if not isinstance(value, str) or not value or Path(value).is_absolute():
            return False
        return (root / value).resolve().is_relative_to(root.resolve())
    for ident, task in by_id.items():
        state = task.get("state")
        if state not in STATES:
            errors.append(f"{ident}: invalid state")
        deps = task.get("depends_on")
        if not isinstance(deps, list) or any(not isinstance(d, str) for d in deps):
            errors.append(f"{ident}: depends_on must be string list")
            deps = []
        if len(deps) != len(set(deps)):
            errors.append(f"{ident}: duplicate dependency")
        for dep in deps:
            if dep not in by_id:
                errors.append(f"{ident}: unknown dependency {dep}")
            elif state in {"ready", "in_progress"} and by_id[dep].get("state") != "complete":
                # The coordinator gate reviews several implemented leaves together.
                if not (ident == "native-leaves-integrate" and by_id[dep].get("state") == "implemented"):
                    errors.append(f"{ident}: prerequisite {dep} is not complete")
        for key in ("allowed_paths", "read_only_paths"):
            values = task.get(key)
            if not isinstance(values, list) or any(not safe_path(v) for v in values):
                errors.append(f"{ident}: invalid {key}")
        for key in ("contract", "task"):
            value = task.get(key)
            if value is not None and (not safe_path(value) or not (root / value).is_file()):
                errors.append(f"{ident}: missing or unsafe {key}")
        if state != "planned":
            if not task.get("task") or not task.get("allowed_paths") or not task.get("checks"):
                errors.append(f"{ident}: specified task requires task file, owned paths and checks")
            if task.get("owner") != "coordinator" and not task.get("contract"):
                errors.append(f"{ident}: implementation task needs a contract")
        if state == "planned" and not task.get("reason"):
            errors.append(f"{ident}: planned task needs a reason")
        contract = task.get("contract")
        if contract and safe_path(contract) and (root / contract).is_file():
            rev = task.get("contract_revision")
            content = (root / contract).read_text()
            if not isinstance(rev, int) or rev < 1 or f"Revision: {rev} ·" not in content:
                errors.append(f"{ident}: contract revision mismatch")
            if state in {"ready", "in_progress"} and "State: specified" not in content:
                errors.append(f"{ident}: contract is not specified")
        if state in {"implemented", "complete"}:
            if not (root / "work" / "handoffs" / (ident + ".md")).is_file():
                errors.append(f"{ident}: missing handoff evidence")
    visiting, visited = set(), set()
    def visit(ident):
        if ident in visiting:
            errors.append(f"dependency cycle at {ident}")
            return
        if ident in visited:
            return
        visiting.add(ident)
        deps = by_id[ident].get("depends_on", [])
        if isinstance(deps, list):
            for dep in deps:
                if isinstance(dep, str) and dep in by_id:
                    visit(dep)
        visiting.remove(ident)
        visited.add(ident)
    for ident in by_id:
        visit(ident)
    return errors

def main():
    root = Path(__file__).resolve().parents[2]
    try:
        document = json.loads((root / "work/manifest.json").read_text())
        problems = validate(root, document)
    except (OSError, ValueError, TypeError, AttributeError) as exc:
        print(f"Manifest could not be validated: {exc}", file=sys.stderr)
        return 2
    if problems:
        print("\n".join(problems), file=sys.stderr)
        return 1
    ready = [t["id"] for t in document["tasks"] if t["state"] == "ready"]
    print(f"Manifest valid: {len(document['tasks'])} tasks. Ready: {', '.join(ready) or 'none'}.")
    print("Metadata only; source imports, implementation and review evidence are not verified.")
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
