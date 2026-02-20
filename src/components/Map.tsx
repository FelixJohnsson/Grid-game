import * as T from "../api/types";
import Person from "./Person";
import { useMemo, useState } from "react";

type Props = {
  world: T.CleanedTile[][] | undefined;
};

const Map = ({ world }: Props) => {
  const [tileSize, setTileSize] = useState(9);
  const [showGrid, setShowGrid] = useState(false);
  const [tooltip, setTooltip] = useState<{
    text: string;
    x: number;
    y: number;
  } | null>(null);

  const claimedBases = useMemo(() => {
    const lookup = new globalThis.Map<string, string[]>();
    if (!world) {
      return lookup;
    }

    for (const row of world) {
      for (const tile of row) {
        const tileEntity = tile.Entity ?? tile.Person;
        if (!tileEntity?.BaseClaimed || !tileEntity.BaseLocation) {
          continue;
        }

        const key = `${tileEntity.BaseLocation.X},${tileEntity.BaseLocation.Y}`;
        const owners = lookup.get(key) ?? [];
        if (!owners.includes(tileEntity.FullName)) {
          owners.push(tileEntity.FullName);
          lookup.set(key, owners);
        }
      }
    }

    return lookup;
  }, [world]);

  const mapStats = useMemo(() => {
    if (!world) {
      return { entities: 0, plants: 0, items: 0, shelters: 0, bases: 0 };
    }

    let entities = 0;
    let plants = 0;
    let items = 0;
    let shelters = 0;
    let bases = 0;

    for (const row of world) {
      for (const tile of row) {
        const tileEntity = tile.Entity ?? tile.Person;
        if (tileEntity) {
          entities++;
        }
        if (tile.Plant) {
          plants++;
        }
        if (tile.Shelter) {
          shelters++;
        }
        if (tile.Items?.length) {
          items += tile.Items.length;
        }
      }
    }

    bases = claimedBases.size;

    return { entities, plants, items, shelters, bases };
  }, [world, claimedBases]);

  const handleMouseEnter = (event: React.MouseEvent, text: string) => {
    const { clientX, clientY } = event;
    setTooltip({ text, x: clientX, y: clientY });
  };

  const handleMouseMove = (event: React.MouseEvent) => {
    const { clientX, clientY } = event;
    setTooltip((current) => (current ? { ...current, x: clientX, y: clientY } : current));
  };

  const handleMouseLeave = () => {
    setTooltip(null);
  };

  const clampTileSize = (nextValue: number) => {
    return Math.max(5, Math.min(16, nextValue));
  };

  const tileTypeLabel = (type: T.TileType): string => {
    if (type === T.TileType.Grass) {
      return "Grass";
    }
    if (type === T.TileType.Water) {
      return "Water";
    }
    return "Mountain";
  };

  const tileBaseClass = (type: T.TileType): string => {
    if (type === T.TileType.Grass) {
      return "bg-emerald-500";
    }
    if (type === T.TileType.Water) {
      return "bg-sky-500";
    }
    return "bg-slate-500";
  };

  const plantMarkerClass = (name: string): string => {
    const normalized = name.toLowerCase();
    if (normalized.includes("apple")) {
      return "bg-lime-300 text-lime-900";
    }
    if (normalized.includes("oak")) {
      return "bg-amber-700 text-amber-100";
    }
    if (normalized.includes("grass")) {
      return "bg-emerald-700 text-emerald-100";
    }
    return "bg-orange-700 text-orange-100";
  };

  const buildTooltipText = (
    tile: T.CleanedTile,
    x: number,
    y: number,
    baseOwners: string[]
  ): string => {
    const tileEntity = tile.Entity ?? tile.Person;
    const lines: string[] = [];
    const stickCount =
      tile.Items?.filter((item) => item.Name.toLowerCase().includes("stick")).length ?? 0;
    const stoneCount =
      tile.Items?.filter((item) => item.Name.toLowerCase() === "stone").length ?? 0;

    lines.push(`Location: (x:${x}, y:${y})`);
    lines.push(`Tile: ${tileTypeLabel(tile.Type)}`);

    if (tileEntity) {
      lines.push(`Entity: ${tileEntity.FullName}`);
      lines.push(`Task: ${tileEntity.CurrentTask?.Action ?? "Unknown"}`);
      if (tileEntity.Personalities?.length) {
        lines.push(`Personality: ${tileEntity.Personalities.join(", ")}`);
      }
      if (tileEntity.LastDialogue) {
        lines.push(
          `Last dialogue: "${tileEntity.LastDialogue}"` +
          (tileEntity.DialogueWith ? ` (with ${tileEntity.DialogueWith})` : "")
        );
      }
    }

    if (tile.Plant) {
      const fruitCount = tile.Plant.Fruit?.length ?? 0;
      lines.push(`Plant: ${tile.Plant.Name}`);
      if (fruitCount > 0) {
        lines.push(`Fruit: ${fruitCount}`);
      }
    }

    if (tile.Shelter) {
      lines.push("Shelter: Present");
    }

    if (baseOwners.length > 0) {
      lines.push(`Claimed base: ${baseOwners.join(", ")}`);
    }

    if (tile.Items?.length) {
      lines.push(`Items: ${tile.Items.map((item) => item.Name).join(", ")}`);
      if (stickCount > 0) {
        lines.push(`Sticks on ground: ${stickCount}`);
      }
      if (stoneCount > 0) {
        lines.push(`Stones on ground: ${stoneCount}`);
      }
    }

    return lines.join("\n");
  };

  return (
    <div className="w-full pr-[25rem] pt-3 px-3 relative">
      <div className="border rounded-md bg-white shadow-sm p-3">
        <div className="flex items-center justify-between mb-2">
          <div>
            <h2 className="text-sm font-semibold">World Map</h2>
            <p className="text-xs text-slate-600">
              Entities: {mapStats.entities} | Plants: {mapStats.plants} | Items: {mapStats.items} | Shelters: {mapStats.shelters} | Bases: {mapStats.bases}
            </p>
          </div>
          <div className="flex items-center gap-2 text-xs">
            <button
              type="button"
              className="px-2 py-1 border rounded hover:bg-slate-50"
              onClick={() => setTileSize((value) => clampTileSize(value - 1))}
            >
              -
            </button>
            <span className="min-w-[4.5rem] text-center">Tile {tileSize}px</span>
            <button
              type="button"
              className="px-2 py-1 border rounded hover:bg-slate-50"
              onClick={() => setTileSize((value) => clampTileSize(value + 1))}
            >
              +
            </button>
            <button
              type="button"
              className={`px-2 py-1 border rounded ${
                showGrid ? "bg-slate-900 text-white border-slate-900" : "hover:bg-slate-50"
              }`}
              onClick={() => setShowGrid((value) => !value)}
            >
              Grid
            </button>
          </div>
        </div>

        <div className="border rounded-md overflow-auto bg-slate-100 max-h-[calc(100vh-10rem)]">
          {world ? (
            <div>
              {world.map((row, y) => (
                <div key={`row-${y}`} className="flex">
                  {row.map((tile, x) => {
                    const tileEntity = tile.Entity ?? tile.Person;
                    const baseOwners = claimedBases.get(`${x},${y}`) ?? [];
                    const tooltipText = buildTooltipText(tile, x, y, baseOwners);
                    const stickCount =
                      tile.Items?.filter((item) =>
                        item.Name.toLowerCase().includes("stick")
                      ).length ?? 0;
                    const stoneCount =
                      tile.Items?.filter((item) => item.Name.toLowerCase() === "stone")
                        .length ?? 0;

                    return (
                      <div
                        key={`${y}-${x}`}
                        className={`relative ${tileBaseClass(tile.Type)}`}
                        style={{
                          width: `${tileSize}px`,
                          height: `${tileSize}px`,
                          boxSizing: "border-box",
                          border: showGrid ? "1px solid rgba(15,23,42,0.15)" : "none",
                        }}
                        onMouseEnter={(event) => handleMouseEnter(event, tooltipText)}
                        onMouseMove={handleMouseMove}
                        onMouseLeave={handleMouseLeave}
                      >
                        {tile.Plant ? (
                          <div
                            className={`absolute top-0 left-0 text-[8px] font-bold px-[1px] leading-none ${plantMarkerClass(
                              String(tile.Plant.Name)
                            )}`}
                          >
                            {String(tile.Plant.Name).charAt(0)}
                          </div>
                        ) : null}

                        {tile.Shelter ? (
                          <div className="absolute bottom-0 left-0 text-[8px] font-bold text-yellow-900 bg-yellow-300 px-[1px] leading-none">
                            S
                          </div>
                        ) : null}

                        {tile.Items?.length ? (
                          <div className="absolute bottom-0 right-0 text-[8px] font-bold text-black bg-amber-200 px-[1px] leading-none">
                            I
                          </div>
                        ) : null}

                        {baseOwners.length > 0 ? (
                          <div className="absolute top-0 right-0 text-[8px] font-bold text-sky-900 bg-cyan-200 px-[1px] leading-none">
                            B
                          </div>
                        ) : null}

                        {stickCount > 0 ? (
                          <div className="absolute bottom-0 right-[9px] text-[7px] font-bold text-amber-100 bg-amber-800 px-[1px] leading-none">
                            St
                          </div>
                        ) : null}

                        {stoneCount > 0 ? (
                          <div className="absolute bottom-0 right-[18px] text-[7px] font-bold text-slate-100 bg-slate-700 px-[1px] leading-none">
                            R
                          </div>
                        ) : null}

                        {tileEntity ? (
                          <div className="absolute inset-0 flex justify-center items-center z-10">
                            <Person person={tileEntity} tileSize={tileSize} />
                          </div>
                        ) : null}
                      </div>
                    );
                  })}
                </div>
              ))}
            </div>
          ) : (
            <div className="p-4 text-xs text-slate-600">Loading world map...</div>
          )}
        </div>

        <div className="mt-2 text-xs text-slate-700 flex gap-3 flex-wrap">
          <span><span className="inline-block w-2 h-2 bg-emerald-500 mr-1" />Grass</span>
          <span><span className="inline-block w-2 h-2 bg-sky-500 mr-1" />Water</span>
          <span><span className="inline-block w-2 h-2 bg-slate-500 mr-1" />Mountain</span>
          <span><span className="inline-block w-2 h-2 bg-rose-600 mr-1" />Entity</span>
          <span><span className="inline-block w-2 h-2 bg-lime-300 mr-1" />Plant marker</span>
          <span><span className="inline-block w-2 h-2 bg-cyan-200 mr-1" />Claimed base marker</span>
          <span><span className="inline-block w-2 h-2 bg-yellow-300 mr-1" />Shelter marker</span>
          <span><span className="inline-block w-2 h-2 bg-amber-200 mr-1" />Item marker</span>
          <span><span className="inline-block w-2 h-2 bg-amber-800 mr-1" />Sticks on ground marker</span>
          <span><span className="inline-block w-2 h-2 bg-slate-700 mr-1" />Stone on ground marker</span>
        </div>
      </div>

      {tooltip && (
        <div
          className="fixed bg-slate-900 text-white text-xs rounded px-2 py-1 z-[60] shadow-md whitespace-pre-line max-w-sm pointer-events-none"
          style={{ top: tooltip.y + 14, left: tooltip.x + 14 }}
        >
          {tooltip.text}
        </div>
      )}
    </div>
  );
};

export default Map;
