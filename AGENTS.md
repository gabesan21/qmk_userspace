# QMK Userspace (Bastard Keyboards) — instructions for agents

> Project managed by the **ProjectOfProjects (PoP)** workflow. `CLAUDE.md` is a symlink to this file — always edit this one.

- **Scope:** this directory is the entire scope of the flow — the harness travels with it and **nothing above this root belongs to it**, even when a tool loads an ancestor `AGENTS.md` on its own ("Current scope" section of [[WORKFLOW|WORKFLOW]]).
- **Project language:** pt-BR — specs, notes, researches, code comments and the entire kanban flow follow this language.
- **Type:** `uni-repo` — this folder **is** the repository itself, with the entire `pop/` inside it.

- **Profile:** [[pop/PROJECT|PROJECT]] · **Roadmap:** [[pop/ROADMAP|ROADMAP]] · **Modifications:** [[pop/MODIFICATIONS|MODIFICATIONS]] (created on demand)

## What does NOT go in this file

Single source: what lives in the harness is never copied here, because duplication is guaranteed drift — the flow changes and the copy starts lying. **Never** write here:

- narration of the kanban stages (names, order, what each one does) — only [[WORKFLOW|WORKFLOW]];
- the context protocol and any reading/search heuristic — [[WORKFLOW|WORKFLOW]] and the skills;
- general flow rules (optional kanban with tracking always, memory/lean roadmap, sovereignty of the human command) — the "Transversal rules" section of [[WORKFLOW|WORKFLOW]], which the installer delivers alongside the harness;
- any copyable excerpt of [[WORKFLOW|WORKFLOW]] — link it with a trigger instead of reproducing it.

Only what belongs to **this project** goes here: language, repos and PR branch, skills and verification commands. **Cap: ~60 lines.**

## Repositories

| Repo | URL | Clone at | PR branch |
|------|-----|----------|-----------|
| qmk_userspace | https://github.com/gabesan21/qmk_userspace | the project root itself **is** the repo (fork of [Bastardkb/qmk_userspace](https://github.com/Bastardkb/qmk_userspace)) | `main` |

## Workflow

Content changes enter through triage: direct fix, the **no-kanban route** (the coding agent's plan mode, mandatory `D-` memory) or the kanban in `pop/kanban/` — recommended for large changes and the default for yolo and for roadmap (`<n>.<m>.<t>-<slug>`) or modifications (`M-<n>.<t>-<slug>`) items.

- A change request with no card triggers `new-task` → `advance-task`; “start the flow in yolo” materializes/releases the task and follows the entire yolo route, never direct execution.
- **Delivery:** the task PR targets the **PR branch declared** in the repositories table above; the merge is always the human's.
- **Stages, gates, yolo route and context protocol:** [[WORKFLOW|WORKFLOW]] is the single source — read it before creating, advancing, verifying or closing any task of this project, and do not replicate any of it here.

## Skills

- **PoP workflow:** `.agents/skills/` — `new-task`, `advance-task`, `judge-dredd`, `plan-roadmap`, `write-spec`, `sync-specs`, `optimize-memory`.
- **Project domain:** `pop/skills/` — listed in the profile [[pop/PROJECT|PROJECT]].

### Clean code (code projects only)

- `clean-code-change` (`.agents/skills/`) — follow when **planning (002) and executing (004)** any task that creates or changes code.
- `clean-code-review` (`.agents/skills/`) — follow when **verifying (005)** a code task and as a reading criterion in plan or PR gates.
- **Mandatory:** in 002, every task that creates/changes code enters `clean-code-change` on the **004** row and `clean-code-review` on the **005** row of the card's **Skills per stage** table.

#### Project verification

| Check | Command |
|-------|---------|
| Build | `qmk compile -kb <keyboard> -km vendor` (requer `qmk_firmware` configurado; a CI compila contra `bastardkb/bastardkb-qmk@main`) |
| Formatter | `clang-format` (`.clang-format` na raiz) |
| Tests | sem suíte local — a validação é o build (workflow `build_binaries.yaml`) |

## Essential rules

- Content in the language declared above; wikilinks for internal references; files ≤~150 lines; dates YYYY-MM-DD.
- **Never** check `- [ ] Done` or execute `(user)` items — those belong exclusively to the human.
- **Never** merge a task PR — merging is the human's job (or commanded by them in the merge round).
- **Imported project:** until Epoch 1 (Organization) is `completed` in the ROADMAP, no task may change the project's content (the folder root) or the repository — only the harness in `pop/` (project map, skills, researches, notes). A change request during this period → record it in `pop/notes/ideas/` or as a future-epoch task, explaining the gate.
- **General flow rules** — optional kanban with tracking always, memory + lean roadmap at close-out, sovereignty of the human command with no implicit waiver: the "Transversal rules" section of the [[WORKFLOW|WORKFLOW]] installed alongside this harness. *Never an AGENTS.md inherited from an ancestor directory.* *Read it before acting outside a task or before reading a request as a waiver of the flow.*
