import * as T from "../api/types";

type Props = {
  person: T.PersonCleaned;
  tileSize: number;
};

const speechBubbleDurationMs = 4000;

const Person = ({ person, tileSize }: Props) => {
  const markerSize = Math.max(4, Math.floor(tileSize * 0.8));
  const showInitial = tileSize >= 10;
  const markerFontSize = Math.max(6, Math.min(9, Math.floor(tileSize * 0.5)));
  const speechFontSize = Math.max(8, Math.min(10, Math.floor(tileSize * 0.85)));
  const speechMaxWidth = Math.max(90, tileSize * 12);

  const hasRecentSpokenDialogue =
    person.DialogueMode === "Said" &&
    !!person.LastDialogue &&
    !!person.DialogueAtMs &&
    Date.now() - person.DialogueAtMs <= speechBubbleDurationMs;

  return (
    <div key={person.FullName} className="relative flex items-center justify-center pointer-events-none">
      {hasRecentSpokenDialogue ? (
        <div
          className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 z-20 bg-white text-slate-900 border border-slate-300 rounded px-1 py-[1px] shadow-sm whitespace-normal text-center leading-tight"
          style={{ maxWidth: `${speechMaxWidth}px`, fontSize: `${speechFontSize}px` }}
        >
          {person.LastDialogue}
          <div className="absolute -bottom-[3px] left-1/2 -translate-x-1/2 w-1.5 h-1.5 bg-white border-r border-b border-slate-300 rotate-45" />
        </div>
      ) : null}

      <div
        className="bg-rose-600 border border-rose-900 rounded-full flex justify-center items-center text-white font-semibold shadow-sm"
        style={{
          width: markerSize,
          height: markerSize,
          fontSize: `${markerFontSize}px`,
          lineHeight: `${markerFontSize}px`,
        }}
      >
        {showInitial ? person.FullName[0] : ""}
      </div>
    </div>
  );
};

export default Person;
