import * as T from "../api/types";
import Person from "./Person";
import { useState } from "react";

type Props = {
  world: T.CleanedTile[][] | undefined;
  grab: (item: T.Item, person: T.Person) => void;
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
                    width: "30px",
                    height: "30px",
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
                          className={
                            tile.Persons && tile.Persons.length > 0
                              ? "hidden"
                              : ""
                          }
                          onMouseEnter={(e) => handleMouseEnter(e, item.Name)}
                          onMouseLeave={handleMouseLeave}
                        >
                          {item.Name[0]}
                        </div>
                      ))}
                    </div>
                  ) : null}

                  {tile.Persons ? (
                    <div>
                      {tile.Persons.map((person, index) => (
                        <Person
                          key={index}
                          person={person}
                          currentTile={tile}
                          grab={grab}
                          onMouseEnter={handleMouseEnter}
                          onMouseLeave={handleMouseLeave}
                        />
                      ))}
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
          style={{ top: tooltip.y - 300, left: tooltip.x - 50 }}
        >
          {tooltip.text}
        </div>
      )}
    </div>
  );
};

export default Map;
