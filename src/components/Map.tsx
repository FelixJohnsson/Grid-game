import * as T from "../api/types";
import Building from "./Building";
import Person from "./Person";

type Props = {
  world: T.World["tiles"] | undefined;
  grab: (item: T.Item, person: T.Person) => void;
};

const Map = ({ world, grab }: Props) => {
  return (
    <div className="w-full flex justify-center pt-6">
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
                      tile.type === T.TileType.Grass
                        ? "green"
                        : tile.type === T.TileType.Water
                        ? "blue"
                        : "gray",
                  }}
                >
                  {tile.building ? <Building building={tile.building} /> : null}

                  {tile.items ? (
                    <div>
                      {tile.items.map((item, index) => (
                        <div
                          key={index}
                          className={
                            tile.persons && tile.persons.length > 0
                              ? "hidden"
                              : ""
                          }
                        >
                          {item.Name[0]}
                        </div>
                      ))}
                    </div>
                  ) : null}

                  {tile.persons ? (
                    <div>
                      {tile.persons.map((person, index) => (
                        <Person
                          key={index}
                          person={person}
                          currentTile={tile}
                          grab={grab}
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
    </div>
  );
};

export default Map;
