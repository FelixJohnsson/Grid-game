import * as T from "../api/types";
import Person from "./Person";
import { useState } from "react";

type Props = {
  world: T.CleanedTile[][] | undefined;
  grab: (item: T.Item, person: T.PersonCleaned) => void;
};

const Map = ({ world, grab }: Props) => {
  const [tooltip, setTooltip] = useState<{
    text: string;
    x: number;
    y: number;
  } | null>(null);

  const handleMouseEnter = (event: React.MouseEvent, text: string) => {
    const { clientX, clientY } = event;
    setTooltip({ text, x: clientX, y: clientY });
  };

  const handleMouseLeave = () => {
    setTooltip(null);
  };

  return (
    <div className="w-full flex justify-center pt-6 relative">
      {world ? (
        <div>
          {world.map((row, y) => (
            <div key={y} className="flex">
              {row.map((tile, x) => (
                <div
                  key={x}
                  style={{
                    width: "15px",
                    height: "15px",
                    position: "relative",
                    backgroundColor:
                      tile.Type === T.TileType.Grass
                        ? "green"
                        : tile.Type === T.TileType.Water
                        ? "blue"
                        : "gray",
                  }}
                >
                  {tile.Items ? (
                    <div>
                      {tile.Items.map((item, index) => (
                        <div
                          key={index}
                          className={tile.Person ? "hidden" : ""}
                          onMouseEnter={(e) => handleMouseEnter(e, item.Name)}
                          onMouseLeave={handleMouseLeave}
                        >
                          <p className="text-xs">{item.Name[0]}</p>
                        </div>
                      ))}
                    </div>
                  ) : null}

                  {tile.Plant ? (
                    <div>
                      <div
                        className="bg-orange-800"
                        onMouseEnter={(e) =>
                          handleMouseEnter(e, tile.Plant.Name)
                        }
                        onMouseLeave={handleMouseLeave}
                      >
                        <p className="text-xs">{tile.Plant.Name[0]}</p>
                      </div>
                    </div>
                  ) : null}

                  {tile.Shelter ? (
                    <div>
                      <div
                        className="bg-yellow-800"
                        onMouseEnter={(e) => handleMouseEnter(e, "Shelter")}
                        onMouseLeave={handleMouseLeave}
                      >
                        <p className="text-xs">Shelter</p>
                      </div>
                    </div>
                  ) : null}

                  {tile.Person ? (
                    <div>
                      <Person
                        person={tile.Person}
                        currentTile={tile}
                        grab={grab}
                        onMouseEnter={handleMouseEnter}
                        onMouseLeave={handleMouseLeave}
                      />
                    </div>
                  ) : null}
                </div>
              ))}
            </div>
          ))}
        </div>
      ) : null}

      {tooltip && (
        <div
          className="absolute bg-gray-700 text-white text-xs rounded p-1"
          style={{ top: tooltip.y, left: tooltip.x }}
        >
          {tooltip.text}
        </div>
      )}
    </div>
  );
};

export default Map;
