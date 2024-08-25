import * as T from "../api/types";

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
              <div key={y} style={{ display: "flex" }}>
                {row.map((tile, x) => (
                  <div
                    key={x}
                    style={{
                      width: "20px",
                      height: "20px",
                      backgroundColor:
                        tile.type === T.TileType.Grass
                          ? "green"
                          : tile.type === T.TileType.Water
                          ? "blue"
                          : "gray",
                    }}
                  >
                    {tile.building ? (
                      <div className={tile.building.Color}>
                        {tile.building.Icon}
                      </div>
                    ) : null}
                    {tile.persons ? (
                      <div>
                        {tile.persons.map((person, i) => (
                          <div key={i} className="bg-orange-500">
                            {person.Initials}
                          </div>
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
