import copy
from pathlib import Path
import tempfile
import unittest
from verify_manifest import validate

class ManifestTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        (self.root / "contract.md").write_text("Revision: 1 · State: specified")
        (self.root / "task.md").write_text("task")
        self.doc = {"schema_version": 1, "tasks": [{
            "id": "leaf", "state": "ready", "depends_on": [],
            "contract": "contract.md", "contract_revision": 1, "task": "task.md",
            "owner": None, "allowed_paths": ["src/leaf/"],
            "read_only_paths": ["contract.md"], "checks": ["test-command"]}]}
    def test_valid_ready_task(self):
        self.assertEqual(validate(self.root, self.doc), [])
    def test_missing_task(self):
        (self.root / "task.md").unlink()
        self.assertTrue(any("missing or unsafe task" in e for e in validate(self.root, self.doc)))
    def test_revision_drift(self):
        self.doc["tasks"][0]["contract_revision"] = 2
        self.assertTrue(any("revision mismatch" in e for e in validate(self.root, self.doc)))
    def test_unknown_dependency(self):
        self.doc["tasks"][0]["depends_on"] = ["missing"]
        self.assertTrue(any("unknown dependency" in e for e in validate(self.root, self.doc)))
    def test_cycle(self):
        self.doc["tasks"][0]["depends_on"] = ["leaf"]
        self.assertTrue(any("cycle" in e for e in validate(self.root, self.doc)))
    def test_incomplete_prerequisite(self):
        second = copy.deepcopy(self.doc["tasks"][0]); second["id"] = "second"; second["depends_on"] = ["leaf"]
        self.doc["tasks"].append(second)
        self.assertTrue(any("not complete" in e for e in validate(self.root, self.doc)))
    def test_escape(self):
        self.doc["tasks"][0]["allowed_paths"] = ["../outside"]
        self.assertTrue(any("invalid allowed_paths" in e for e in validate(self.root, self.doc)))
    def test_missing_evidence(self):
        self.doc["tasks"][0]["state"] = "complete"
        self.assertTrue(any("missing handoff" in e for e in validate(self.root, self.doc)))
    def test_duplicate(self):
        self.doc["tasks"].append(copy.deepcopy(self.doc["tasks"][0]))
        self.assertTrue(any("duplicate task" in e for e in validate(self.root, self.doc)))
    def test_planned_cannot_become_ready_without_spec(self):
        self.doc["tasks"][0].update(contract=None, task=None, checks=[])
        self.assertTrue(validate(self.root, self.doc))

if __name__ == "__main__":
    unittest.main()
