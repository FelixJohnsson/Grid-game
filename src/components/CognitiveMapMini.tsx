import { useMemo } from "react";
import * as T from "../api/types";

type Props = {
  knownTiles: T.CognitiveMapKnownTile[];
  center: T.Location;
  isLoading?: boolean;
  error?: string | null;
};

const VIEW_SIZE = 31;
const TILE_SIZE_PX = 5;

const keyForLocation = (location: T.Location): string => {
  return `${location.X},${location.Y}`;
};

const getTileClassName = (tile: T.CognitiveMapKnownTile | undefined): string => {
  if (!tile) {
    return "bg-slate-900";
  }

  if (tile.Entity?.IsAlive) {
    const species = tile.Entity.SpeciesType?.toLowerCase() ?? "";
    if (species.includes("homo sapiens")) {
      return "bg-red-500";
    }
    if (species.includes("canis lupus")) {
      return "bg-gray-400";
    }
    return "bg-purple-500";
  }

  if (tile.Plant?.IsAlive) {
    const plantName = tile.Plant.Name?.toLowerCase() ?? "";
    if (plantName.includes("apple")) {
      return "bg-lime-300";
    }
    if (plantName.includes("oak")) {
      return "bg-amber-700";
    }
    if (plantName.includes("grass")) {
      return "bg-emerald-700";
    }
    if (plantName.includes("flower")) {
      return "bg-pink-500";
    }
    return "bg-emerald-500";
  }

  if (tile.TileType === T.TileType.Water) {
    return "bg-blue-500";
  }
  if (tile.TileType === T.TileType.Grass) {
    return "bg-green-500";
  }

  return "bg-gray-500";
};

const CognitiveMapMini = ({ knownTiles, center, isLoading, error }: Props) => {
  const knownTileCount = Array.isArray(knownTiles) ? knownTiles.length : 0;

  const knownTileLookup = useMemo(() => {
    const lookup = new Map<string, T.CognitiveMapKnownTile>();
    const tiles = Array.isArray(knownTiles) ? knownTiles : [];
    tiles.forEach((tile) => {
      lookup.set(keyForLocation(tile.Location), tile);
    });
    return lookup;
  }, [knownTiles]);

  const minX = center.X - Math.floor(VIEW_SIZE / 2);
  const minY = center.Y - Math.floor(VIEW_SIZE / 2);

  const cells: JSX.Element[] = [];
  for (let y = 0; y < VIEW_SIZE; y++) {
    for (let x = 0; x < VIEW_SIZE; x++) {
      const worldX = minX + x;
      const worldY = minY + y;
      const knownTile = knownTileLookup.get(`${worldX},${worldY}`);
      const isCenter = worldX === center.X && worldY === center.Y;

      let title = `Unknown (${worldX}, ${worldY})`;
      if (knownTile) {
        title = `Known (${worldX}, ${worldY})`;
        if (knownTile.Entity?.FullName) {
          title += ` - Entity: ${knownTile.Entity.FullName}`;
        } else if (knownTile.Plant?.Name) {
          title += ` - Plant: ${knownTile.Plant.Name}`;
        }
      }
      if (isCenter) {
        title += " - Current position";
      }

      cells.push(
        <div
          key={`${worldX}-${worldY}`}
          className={`${getTileClassName(knownTile)} ${isCenter ? "ring-1 ring-white" : ""}`}
          style={{ width: TILE_SIZE_PX, height: TILE_SIZE_PX, border: "0.5px solid rgba(0,0,0,0.2)" }}
          title={title}
        />
      );
    }
  }

  return (
    <div className="border rounded p-2 space-y-1">
      <p className="font-semibold">Cognitive map (mini)</p>
      <p>
        Known tiles: {knownTileCount} | Center: ({center.X}, {center.Y})
      </p>
      {error ? <p className="text-red-700">{error}</p> : null}
      {isLoading && knownTileCount === 0 ? <p>Loading cognitive map...</p> : null}

      <div
        className="inline-grid bg-black p-1 rounded"
        style={{ gridTemplateColumns: `repeat(${VIEW_SIZE}, ${TILE_SIZE_PX}px)` }}
      >
        {cells}
      </div>

      <div className="flex flex-wrap gap-2">
        <span className="inline-flex items-center gap-1">
          <span className="inline-block w-2 h-2 bg-red-500" /> human
        </span>
        <span className="inline-flex items-center gap-1">
          <span className="inline-block w-2 h-2 bg-gray-400" /> wolf
        </span>
        <span className="inline-flex items-center gap-1">
          <span className="inline-block w-2 h-2 bg-blue-500" /> water
        </span>
        <span className="inline-flex items-center gap-1">
          <span className="inline-block w-2 h-2 bg-slate-900" /> unknown
        </span>
      </div>
    </div>
  );
};

export default CognitiveMapMini;
