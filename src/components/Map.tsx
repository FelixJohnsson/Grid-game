import * as T from "../api/types";
import Building from "./Building";
import Person from "./Person";

type Props = {
  world: T.World | undefined;
};

const Map = ({ world }: Props) => {
  return (
    <div className="w-full flex justify-center pt-6">
      {world ? (
        <div>
          <div>
            {world.tiles.map((row, y) => (
              <div key={y} className="flex">
                {row.map((tile, x) => (
                  <div
                    key={x}
                    style={{
                      width: "30px",
                      height: "30px",
                      backgroundColor:
                        tile.type === T.TileType.Grass
                          ? "green"
                          : tile.type === T.TileType.Water
                          ? "blue"
                          : "gray",
                    }}
                  >
                    {tile.building ? (
                      <Building building={tile.building} />
                    ) : null}
                    {tile.persons ? (
                      <div>
                        {tile.persons.map((person) => (
                          <Person person={person} />
                        ))}
                      </div>
                    ) : null}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
};

export default Map;
