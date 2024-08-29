type Props = {
  move: (direction: "up" | "down" | "left" | "right") => void;
};

const MoveControls = ({ move }: Props) => {
  return (
    <div>
      <h1>Move</h1>
      <button className="border p-4 cursor-pointer" onClick={() => move("up")}>
        Up
      </button>
      <button
        className="border p-4 cursor-pointer"
        onClick={() => move("down")}
      >
        Down
      </button>
      <button
        className="border p-4 cursor-pointer"
        onClick={() => move("left")}
      >
        Left
      </button>
      <button
        className="border p-4 cursor-pointer"
        onClick={() => move("right")}
      >
        Right
      </button>
    </div>
  );
};

export default MoveControls;
