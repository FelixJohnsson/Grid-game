import { useCallback, useEffect, useRef, useState } from "react";
import * as api from "./api/api";
import * as T from "./api/types";
import Map from "./components/Map";
import BehaviorInspector from "./components/BehaviorInspector";

const getTileEntity = (tile: T.CleanedTile): T.PersonCleaned | undefined => {
  return tile.Entity ?? tile.Person;
};

function App() {
  const [persons, setPersons] = useState<T.PersonCleaned[]>([]);
  const [world, setWorld] = useState<T.CleanedTile[][]>();
  const [isResettingWorld, setIsResettingWorld] = useState(false);
  const didWarnMissingEntityShape = useRef(false);

  const applyWorldUpdate = useCallback((data: T.CleanedTile[][]) => {
    setWorld(data);

    const nextPersons: T.PersonCleaned[] = [];
    data.forEach((row) => {
      row.forEach((tile) => {
        const tileEntity = getTileEntity(tile);
        if (tileEntity) {
          nextPersons.push(tileEntity);
        }
      });
    });

    if (
      !didWarnMissingEntityShape.current &&
      nextPersons.length === 0 &&
      data.length > 0
    ) {
      console.warn(
        "No entities found in /world payload. Check backend tile entity keys."
      );
      didWarnMissingEntityShape.current = true;
    }

    setPersons(nextPersons);
  }, []);

  useEffect(() => {
    const continuslyUpdate = () => {
      api.getWorld().then((data) => {
        applyWorldUpdate(data);
      });
    };

    // Fetch the initial world state
    continuslyUpdate();

    // Set up the interval to continuously update the world state
    const intervalId = setInterval(continuslyUpdate, 500);

    // Clear the interval when the component unmounts to prevent memory leaks and infinite loops
    return () => clearInterval(intervalId);
  }, [applyWorldUpdate]);

  const handleResetWorld = async () => {
    if (isResettingWorld) {
      return;
    }

    setIsResettingWorld(true);
    try {
      const data = await api.resetWorld();
      applyWorldUpdate(data);
    } catch (error) {
      console.error("Failed to reset world", error);
    } finally {
      setIsResettingWorld(false);
    }
  };

  return (
    <div className="App relative min-h-screen bg-slate-200">
      <div className="fixed left-2 top-2 z-50">
        <button
          className="bg-slate-900 hover:bg-black text-white text-xs px-3 py-2 rounded shadow disabled:opacity-60 disabled:cursor-not-allowed"
          onClick={handleResetWorld}
          disabled={isResettingWorld}
        >
          {isResettingWorld ? "Resetting world..." : "Reset world"}
        </button>
      </div>
      <div className="fixed right-2 top-2 z-50 w-96 max-h-[calc(100vh-1rem)] overflow-y-auto">
        <BehaviorInspector persons={persons} />
      </div>
      <Map world={world} />
    </div>
  );
}

export default App;
