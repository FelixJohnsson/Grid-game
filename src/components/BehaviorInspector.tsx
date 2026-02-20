import { useEffect, useMemo, useRef, useState } from "react";
import * as T from "../api/types";
import * as api from "../api/api";
import CognitiveMapMini from "./CognitiveMapMini";

type Props = {
  persons: T.PersonCleaned[];
};

type Snapshot = {
  timestamp: number;
  x: number;
  y: number;
  thought: string;
  taskAction: string;
  taskTarget: string;
  taskPriority: number;
  taskActive: boolean;
};

const HISTORY_LIMIT = 240; // ~2 minutes at 500ms polling
const ANALYSIS_WINDOW_MS = 30_000;
const STUCK_THRESHOLD_SECONDS = 10;

const formatTask = (snapshot: Snapshot): string => {
  if (!snapshot.taskTarget) {
    return snapshot.taskAction;
  }
  return `${snapshot.taskAction} -> ${snapshot.taskTarget}`;
};

const isSameTask = (a: Snapshot, b: Snapshot): boolean => {
  return (
    a.taskAction === b.taskAction &&
    a.taskTarget === b.taskTarget &&
    a.taskActive === b.taskActive
  );
};

const formatSeconds = (value: number): string => {
  if (!Number.isFinite(value)) {
    return "N/A";
  }
  return `${Math.round(value)}s`;
};

const BehaviorInspector = ({ persons }: Props) => {
  const [selectedFullName, setSelectedFullName] = useState("");
  const [cognitiveMapTiles, setCognitiveMapTiles] = useState<T.CognitiveMapKnownTile[]>([]);
  const [isCognitiveMapLoading, setIsCognitiveMapLoading] = useState(false);
  const [cognitiveMapError, setCognitiveMapError] = useState<string | null>(null);
  const historyRef = useRef<Record<string, Snapshot[]>>({});
  const lastCognitiveMapErrorRef = useRef<string | null>(null);

  useEffect(() => {
    const activeNames = new Set(persons.map((person) => person.FullName));
    const now = Date.now();

    persons.forEach((person) => {
      const history = historyRef.current[person.FullName] ?? [];
      history.push({
        timestamp: now,
        x: person.Location.X,
        y: person.Location.Y,
        thought: person.Thinking ?? "",
        taskAction: person.CurrentTask?.Action ?? "Unknown",
        taskTarget: person.CurrentTask?.Target ?? "",
        taskPriority: person.CurrentTask?.Priority ?? 0,
        taskActive: person.CurrentTask?.IsActive ?? false,
      });

      if (history.length > HISTORY_LIMIT) {
        history.splice(0, history.length - HISTORY_LIMIT);
      }
      historyRef.current[person.FullName] = history;
    });

    Object.keys(historyRef.current).forEach((name) => {
      if (!activeNames.has(name)) {
        delete historyRef.current[name];
      }
    });

    if (!selectedFullName || !activeNames.has(selectedFullName)) {
      setSelectedFullName(persons[0]?.FullName ?? "");
    }
  }, [persons, selectedFullName]);

  useEffect(() => {
    if (!selectedFullName) {
      setCognitiveMapTiles([]);
      setCognitiveMapError(null);
      setIsCognitiveMapLoading(false);
      lastCognitiveMapErrorRef.current = null;
      return;
    }

    setCognitiveMapTiles([]);
    setCognitiveMapError(null);
    lastCognitiveMapErrorRef.current = null;

    let cancelled = false;

    const loadCognitiveMap = async (silent: boolean) => {
      if (!silent) {
        setIsCognitiveMapLoading(true);
      }
      try {
        const tiles = await api.getEntityCognitiveMap(selectedFullName);
        if (cancelled) {
          return;
        }
        setCognitiveMapTiles(tiles);
        setCognitiveMapError(null);
        lastCognitiveMapErrorRef.current = null;
      } catch (error) {
        if (cancelled) {
          return;
        }
        const message =
          error instanceof Error
            ? error.message
            : "Unable to load cognitive map for this character.";

        if (lastCognitiveMapErrorRef.current !== message) {
          console.error("Failed to load cognitive map", error);
          lastCognitiveMapErrorRef.current = message;
        }

        setCognitiveMapError(message);
      } finally {
        if (!silent && !cancelled) {
          setIsCognitiveMapLoading(false);
        }
      }
    };

    void loadCognitiveMap(false);
    const intervalId = setInterval(() => {
      void loadCognitiveMap(true);
    }, 2000);

    return () => {
      cancelled = true;
      clearInterval(intervalId);
    };
  }, [selectedFullName]);

  const selectedPerson = useMemo(
    () => persons.find((person) => person.FullName === selectedFullName),
    [persons, selectedFullName]
  );

  const history = selectedFullName ? historyRef.current[selectedFullName] ?? [] : [];
  const now = Date.now();
  const recentSnapshots = history.filter(
    (snapshot) => now-snapshot.timestamp <= ANALYSIS_WINDOW_MS
  );

  let movedTiles = 0;
  let movementEvents = 0;
  let taskTransitions = 0;
  let thoughtChanges = 0;

  for (let i = 1; i < recentSnapshots.length; i++) {
    const previous = recentSnapshots[i - 1];
    const current = recentSnapshots[i];

    const stepDistance = Math.abs(current.x - previous.x) + Math.abs(current.y - previous.y);
    if (stepDistance > 0) {
      movementEvents++;
      movedTiles += stepDistance;
    }

    if (!isSameTask(previous, current)) {
      taskTransitions++;
    }

    if (previous.thought !== current.thought) {
      thoughtChanges++;
    }
  }

  let secondsSinceLastMove = Number.POSITIVE_INFINITY;
  for (let i = history.length - 1; i > 0; i--) {
    const current = history[i];
    const previous = history[i - 1];
    if (current.x !== previous.x || current.y !== previous.y) {
      secondsSinceLastMove = (now - current.timestamp) / 1000;
      break;
    }
  }

  const latestSnapshot = history[history.length - 1];
  let taskDurationSeconds = 0;
  if (latestSnapshot) {
    let startTimestamp = latestSnapshot.timestamp;
    for (let i = history.length - 2; i >= 0; i--) {
      if (isSameTask(history[i], latestSnapshot)) {
        startTimestamp = history[i].timestamp;
      } else {
        break;
      }
    }
    taskDurationSeconds = (now - startTimestamp) / 1000;
  }

  const behaviorEvents: { timestamp: number; text: string }[] = [];
  const startIndex = Math.max(1, history.length - 40);
  for (let i = startIndex; i < history.length; i++) {
    const previous = history[i - 1];
    const current = history[i];

    if (current.x !== previous.x || current.y !== previous.y) {
      behaviorEvents.push({
        timestamp: current.timestamp,
        text: `Moved to (${current.x}, ${current.y})`,
      });
    }

    if (!isSameTask(previous, current)) {
      behaviorEvents.push({
        timestamp: current.timestamp,
        text: `Task changed to: ${formatTask(current)}`,
      });
    }

    if (current.thought && current.thought !== previous.thought) {
      behaviorEvents.push({
        timestamp: current.timestamp,
        text: `Thought changed: ${current.thought}`,
      });
    }
  }

  const recentEvents = behaviorEvents.slice(-8).reverse();

  let decisionStatus = "Monitoring";
  let decisionStatusTone = "text-blue-700";
  let decisionSummary = "Collecting behavior data from recent world updates.";

  if (latestSnapshot) {
    if (!latestSnapshot.taskActive && latestSnapshot.taskAction === "Idle") {
      decisionStatus = "Idle";
      decisionStatusTone = "text-gray-700";
      decisionSummary = "No active task is selected right now.";
    } else if (
      latestSnapshot.taskActive &&
      Number.isFinite(secondsSinceLastMove) &&
      secondsSinceLastMove > STUCK_THRESHOLD_SECONDS
    ) {
      decisionStatus = "Potentially stuck";
      decisionStatusTone = "text-red-700";
      decisionSummary =
        "Task is active but movement has not changed recently. The agent may be blocked or looping.";
    } else if (latestSnapshot.taskActive && movedTiles > 0) {
      decisionStatus = "Executing";
      decisionStatusTone = "text-green-700";
      decisionSummary = "Task is active and movement is progressing.";
    } else if (latestSnapshot.taskActive) {
      decisionStatus = "Planning / waiting";
      decisionStatusTone = "text-amber-700";
      decisionSummary = "Task is active but no movement has been observed in the recent window.";
    }
  }

  const decisionSignals: string[] = [];
  if (taskTransitions >= 6) {
    decisionSignals.push("Frequent task switching in the last 30s.");
  }
  if (latestSnapshot?.taskActive && movedTiles === 0) {
    decisionSignals.push("Active task without movement in the last 30s.");
  }
  if (thoughtChanges === 0 && latestSnapshot?.thought) {
    decisionSignals.push("Thought text stayed stable in the recent window.");
  }
  if (selectedPerson?.IsIncapacitated) {
    decisionSignals.push("Entity is incapacitated.");
  }
  if (decisionSignals.length === 0) {
    decisionSignals.push("No obvious anomalies detected in the recent window.");
  }

  return (
    <div className="text-xs bg-white/95 border rounded shadow-lg p-3 space-y-3">
      <h2 className="text-sm font-semibold">Behavior Inspector</h2>

      {persons.length === 0 ? (
        <p>No active persons are currently visible in the world payload.</p>
      ) : (
        <>
          <div>
            <label className="block mb-1 font-medium" htmlFor="inspector-person-select">
              Character
            </label>
            <select
              id="inspector-person-select"
              className="w-full border rounded px-2 py-1"
              value={selectedFullName}
              onChange={(event) => setSelectedFullName(event.target.value)}
            >
              {persons.map((person) => (
                <option key={person.FullName} value={person.FullName}>
                  {person.FullName}
                </option>
              ))}
            </select>
          </div>

          {selectedPerson && latestSnapshot ? (
            <>
              <div className="border rounded p-2 bg-slate-50">
                <p className="font-semibold">Decision snapshot</p>
                <p className={decisionStatusTone}>{decisionStatus}</p>
                <p>{decisionSummary}</p>
              </div>

              <CognitiveMapMini
                knownTiles={cognitiveMapTiles}
                center={selectedPerson.Location}
                isLoading={isCognitiveMapLoading}
                error={cognitiveMapError}
              />

              <div className="border rounded p-2 space-y-1">
                <p className="font-semibold">Current state</p>
                <p>Position: ({selectedPerson.Location.X}, {selectedPerson.Location.Y})</p>
                <p>Task: {formatTask(latestSnapshot)}</p>
                <p>Task priority: {latestSnapshot.taskPriority}</p>
                <p>Task active: {latestSnapshot.taskActive ? "Yes" : "No"}</p>
                <p>Task duration: {formatSeconds(taskDurationSeconds)}</p>
                <p>
                  Thinking: {selectedPerson.Thinking ? selectedPerson.Thinking : "No current thought text"}
                </p>
                <p>
                  Personalities: {selectedPerson.Personalities?.length
                    ? selectedPerson.Personalities.join(", ")
                    : "None"}
                </p>
                <p>
                  Last dialogue: {selectedPerson.LastDialogue
                    ? `"${selectedPerson.LastDialogue}" with ${selectedPerson.DialogueWith ?? "Unknown"}`
                    : "No dialogue yet"}
                </p>
                <p>Incapacitated: {selectedPerson.IsIncapacitated ? "Yes" : "No"}</p>
                <p>Right hand: {selectedPerson.RightArm?.Hand?.Items?.[0]?.Name ?? "Empty"}</p>
                <p>Left hand: {selectedPerson.LeftArm?.Hand?.Items?.[0]?.Name ?? "Empty"}</p>
                <p>Base claimed: {selectedPerson.BaseClaimed ? "Yes" : "No"}</p>
                <p>
                  Base location: {selectedPerson.BaseLocation
                    ? `(${selectedPerson.BaseLocation.X}, ${selectedPerson.BaseLocation.Y})`
                    : "N/A"}
                </p>
                <p>
                  Base stockpile: Food {selectedPerson.BaseStockpile?.Food ?? 0} | Sticks{" "}
                  {selectedPerson.BaseStockpile?.Sticks ?? 0} | Stone{" "}
                  {selectedPerson.BaseStockpile?.Stone ?? 0} | Grass{" "}
                  {selectedPerson.BaseStockpile?.Grass ?? 0}
                </p>
              </div>

              <div className="border rounded p-2 space-y-1">
                <p className="font-semibold">Behavior metrics (last 30s)</p>
                <p>Movement events: {movementEvents}</p>
                <p>Total moved tiles: {movedTiles}</p>
                <p>Task transitions: {taskTransitions}</p>
                <p>Thought changes: {thoughtChanges}</p>
                <p>Seconds since last move: {formatSeconds(secondsSinceLastMove)}</p>
                <p>Relationships tracked: {selectedPerson.Relationships?.length ?? 0}</p>
              </div>

              <div className="border rounded p-2 space-y-1">
                <p className="font-semibold">Decision signals</p>
                {decisionSignals.map((signal) => (
                  <p key={signal}>- {signal}</p>
                ))}
              </div>

              <div className="border rounded p-2 space-y-1">
                <p className="font-semibold">Recent behavior events</p>
                {recentEvents.length === 0 ? (
                  <p>No recent movement, task, or thought changes captured yet.</p>
                ) : (
                  recentEvents.map((event, index) => (
                    <p key={`${event.timestamp}-${index}`}>
                      {new Date(event.timestamp).toLocaleTimeString()}: {event.text}
                    </p>
                  ))
                )}
              </div>
            </>
          ) : (
            <p>Selected character is not available in the latest world snapshot.</p>
          )}
        </>
      )}
    </div>
  );
};

export default BehaviorInspector;
