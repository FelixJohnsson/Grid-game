# Grid Game

Grid-based simulation with autonomous entities, resource gathering, memory-driven behavior, and a React map/inspector UI.

## Quick Start

Run backend and frontend in separate terminals.

### 1) Install frontend dependencies (first time only)

From project root:

```bash
npm install
```

### 2) Start backend API

From project root:

```bash
cd backend
go run .
```

Backend default address: `http://127.0.0.1:8080`

Optional backend display mode (Raylib window):

```bash
cd backend
go run -tags display . -display
```

### 3) Start frontend

From project root (second terminal):

```bash
npm start
```

Frontend default address: `http://localhost:3000`

## Troubleshooting

### `cgo: C compiler "gcc" not found`

- API-only mode: `go run .` (no Raylib display required)
- Raylib display mode: `go run -tags display . -display` (requires a working C toolchain, including `gcc`)

---

## Technical Analysis Report

This report summarizes the current system architecture, what each system is used for, and relative importance.

### Architecture Snapshot

- **Backend** (`Go`): simulation runtime, AI/planning, resource systems, and HTTP API.
- **Frontend** (`React + TypeScript`): live world map, behavior analysis UI, and cognitive-map visualization.
- **Execution model**: continuous simulation loops (brain + motor + plant lifecycle) with frontend polling.

### Runtime Flow (High Level)

1. Backend world is initialized (`tiles`, `entities`, `plants`, `items`).
2. Each entity brain loop evaluates physiology and picks tasks.
3. Motor loop executes movement/pathing toward task targets.
4. Resource and crafting tasks mutate world/base stockpiles.
5. API serializes cleaned world/entity state for frontend polling.
6. Frontend renders map, base/resource markers, and behavior telemetry.

---

## Systems Inventory (Sorted by Importance)

### Critical

| System | Used For | Key Files | Importance |
| --- | --- | --- | --- |
| World state engine | Canonical grid state for tiles/entities/items/plants/shelters; thread-safe mutations and queries | `backend/worldState.go`, `backend/types.go` | Foundation for all simulation behavior |
| Entity brain/planner loop | Need evaluation, task ranking, action dispatch | `backend/brain.go`, `backend/entityActionHandler.go`, `backend/entityBrainFunctions.go` | Core AI decision-making |
| Motor/pathfinding subsystem | Converts task targets into movement over grid | `backend/entityMotorCortex.go`, `backend/AStar.go` | Required for all non-static behaviors |
| HTTP API transport | Exposes world/entity/cognitive-map/reset endpoints to frontend | `backend/server.go`, `backend/cleanedTile.go` | Required for UI integration |
| Frontend world polling + map rendering | Real-time visibility of simulation state | `src/App.tsx`, `src/components/Map.tsx`, `src/components/Person.tsx` | Primary operator/user feedback loop |

### High

| System | Used For | Key Files | Importance |
| --- | --- | --- | --- |
| Cognitive map + memory | Spatial memory, stale-source pruning, map-first resource targeting | `backend/entityVision.go`, `backend/entityMemory.go`, `backend/entityBrainFunctions.go` | Enables non-random, memory-driven behavior |
| Base claim + stockpile logistics | Resource delivery model for shelter/food/stone systems | `backend/entityShelterActions.go`, `backend/types.go` | Backbone for progression and crafting |
| Shelter lifecycle | Claim base, gather sticks, build shelter, post-build fallback | `backend/entityShelterActions.go`, `backend/building.go` | Core survival progression |
| Food hauling pipeline | Collect and deliver food to base with stock target | `backend/entityActions.go`, `backend/entityShelterActions.go`, `backend/entityActionHandler.go` | Sustains longer-term planning |
| Stone system (ground item model) | Stone spawn clusters, detection, memory, hauling, cap | `backend/startWorld.go`, `backend/item.go`, `backend/entityVision.go`, `backend/entityShelterActions.go` | Enables tool-tier progression |
| Crafting recipes engine | Recipe-driven crafting with ingredient checks and stockpile consumption | `backend/craft.go`, `backend/item.go` | Converts gathered resources into utility items |
| Plant lifecycle | Fruit spawn/ripen progression for food ecology | `backend/plantLifecycle.go`, `backend/plant.go` | Maintains renewable food loop |
| Reset world lifecycle | Safe teardown/reseed of simulation loops and entities | `backend/startWorld.go`, `backend/server.go` | Essential for iteration/testing |

### Medium

| System | Used For | Key Files | Importance |
| --- | --- | --- | --- |
| Behavior inspector | Decision snapshots, movement/task metrics, anomaly hints | `src/components/BehaviorInspector.tsx` | High-value debug/analysis tool |
| Cognitive map mini UI | Small per-entity memory visualization | `src/components/CognitiveMapMini.tsx` | Useful for AI observability |
| API type contracts | Shared frontend types for world/cognitive payloads | `src/api/types.ts`, `src/api/api.ts` | Prevents data-shape drift issues |
| Combat/physiology detail systems | Damage, organs, hormones, pain effects | `backend/entityBodyFunctions.go`, `backend/personHostileActions.go`, `backend/brain.go` | Present and functional, but secondary to logistics loop currently |

### Low (current impact)

| System | Used For | Key Files | Importance |
| --- | --- | --- | --- |
| Legacy/unused frontend components | Older UI paths not in active root render flow | `src/components/InformationBar.tsx`, `src/components/MoveControls.tsx`, `src/components/PersonWidget.tsx` | Low runtime impact; maintenance overhead |

---

## What Is Stable Right Now

- Thread-safe world and cognitive-map access is in place.
- Simulation runs with decoupled display mode (`build tags`).
- Base-centric hauling loops exist for sticks, food, and stone.
- Crafting has a recipe model with shelter/ingredient gating.
- Map UI surfaces entities, bases, shelters, and ground resources.
- Backend tests cover major world/vision/shelter/crafting behavior paths.

## Primary Risks / Gaps

- Frontend currently uses polling (no push transport), with fixed intervals.
- Crafting planner is still early-stage (few autonomous craft goals).
- Some action/task enums remain placeholders (`noop`) and are not fully implemented.
- Configuration is mostly hardcoded (addresses, thresholds/constants).
- No persistence layer; state resets on backend restart.

---

## Existing Stone Age Recipes

Implemented in `backend/craft.go`:

- `Stone Knife`
- `Stone Hammer`
- `Stone Axe`
- `Stone Spear`
- `Stone Pickaxe` (requires shelter)
- `Food Box` (requires shelter)
- `Wooden Box` (requires shelter)
- `Wooden Crate` (requires shelter)

---

## Recommended Next Technical Priorities

1. Expand autonomous crafting goals in planner (recipe-aware progression chain).
2. Add shared resource-goal framework for future materials (e.g. clay/flint/obsidian).
3. Add runtime configuration layer for key constants/thresholds.
4. Add frontend tests around map markers and behavior inspector rendering.
5. Add persistence snapshot/load if long-run simulation continuity is required.