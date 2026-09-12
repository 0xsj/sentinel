# <module>: implementation contract

Revision: <integer> · State: draft | specified | superseded
Owner: <coordinator> · Task: <manifest ID>

## Purpose and consumer

Name one responsibility, its first consumer and the reason it needs a separate
module. Distinguish a scheduled consumer from one already implemented.

## Owned files and dependencies

List production/test directories, allowed imports, forbidden imports and shared
integration files owned by someone else. Link prerequisite contract revisions.

## Public API

Give exact names, argument/return types, visibility and error types in the target
language. Resolve every placeholder before marking a task ready. State zero/default
behavior, ownership, mutation, concurrency and lifetime where applicable.

## Observable contract

State input validity, success, refusals, absence, side effects and diagnostic/public
data separately. State ordering, atomicity, cancellation, replay/retry, shutdown
and durability when relevant; explicitly mark irrelevant dimensions.

## Acceptance scenarios

| ID | Given / when | Required observable outcome |
| --- | --- | --- |
| <M01> | <concrete input and operation> | <externally observable result> |
| <M02> | <negative/boundary input> | <result and effects that must not occur> |

Include the most plausible false success. Assign IDs before implementation; adding
coverage is welcome, weakening an existing promise requires coordinator revision.

## Verification

List exact commands from the repository root. Name the checks that do not exist
yet. A filtered command with zero tests is not evidence. Separate pure, real
adapter, native bridge and rendered UI evidence rather than substituting one.

## Exclusions and open decisions

List intentionally absent APIs and product behaviors. Any unresolved question
that changes public behavior makes the task not ready. Implementation mechanics
within the contract are the worker's choice.

## Completion evidence

Link the task handoff, scenario-to-test mapping, command results and review.
Tests written by the implementer are implementation-visible, not a blind oracle.
