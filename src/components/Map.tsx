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
        <div className="flex flex-wrap">
          {world.map((row, y) => (
            <div>
              {row.map((tile, x) => (
                <div
                  key={`${y}-${x}`}
                  className={`relative border ${
                    tile.Type === T.TileType.Grass
                      ? "bg-green-500"
                      : tile.Type === T.TileType.Water
                      ? "bg-blue-500"
                      : "bg-gray-500"
                  }`}
                  style={{ width: "15px", height: "15px" }}
                >
                  {/* Display Items */}
                  {tile.Items &&
                    tile.Items.map((item, index) => (
                      <div
                        key={index}
                        className="absolute inset-0 flex justify-center items-center text-xs z-10"
                        onMouseEnter={(e) => handleMouseEnter(e, item.Name)}
                        onMouseLeave={handleMouseLeave}
                      >
                        <p>{item.Name[0]}</p>
                      </div>
                    ))}

                  {/* Display Plant */}
                  {tile.Plant && (
                    <div
                      className="absolute inset-0 flex justify-center items-center bg-orange-800 text-xs z-20"
                      onMouseEnter={(e) => handleMouseEnter(e, tile.Plant.Name)}
                      onMouseLeave={handleMouseLeave}
                    >
                      <p>{tile.Plant.Name[0]}</p>
                    </div>
                  )}

                  {/* Display Shelter */}
                  {tile.Shelter && (
                    <div
                      className="absolute inset-0 flex justify-center items-center bg-yellow-800 text-xs z-30"
                      onMouseEnter={(e) => handleMouseEnter(e, "Shelter")}
                      onMouseLeave={handleMouseLeave}
                    >
                      <p>S</p>
                    </div>
                  )}

                  {/* Display Person */}
                  {tile.Person && (
                    <div className="absolute inset-0 flex justify-center items-center z-40">
                      <Person
                        person={tile.Person}
                        currentTile={tile}
                        grab={grab}
                        onMouseEnter={handleMouseEnter}
                        onMouseLeave={handleMouseLeave}
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>
          ))}
        </div>
      ) : null}

      {tooltip && (
        <div
          className="absolute bg-gray-700 text-white text-xs rounded p-1"
          style={{ top: tooltip.y - 150, left: tooltip.x }}
        >
          {tooltip.text}
        </div>
      )}
    </div>
  );
};

export default Map;
