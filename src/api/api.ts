import * as T from "./types";

const address = "http://localhost:8080/";

const api = {
  get: async (path: string) => {
    const response = await fetch(address + path);
    return await response.json();
  },
  post: async (path: string, data: any) => {
    const response = await fetch(address + path, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    return await response.json();
  },
};

export const getWorld = async (): Promise<T.CleanedTile[][]> => {
  const data: T.WorldResponse = await api.get("world");
  return data.message;
};

export const movePerson = async (full_name: string, direction: string) => {
  const data = await api.post("move", { full_name, direction });
  return data.message;
};

export const grabItem = async (item: T.Item, person: T.PersonCleaned) => {
  const requestData = {
    ItemName: item.Name,
    FullName: person.FullName,
  };
  const data = await api.post("entityGrab", requestData);
  return data.message;
};

export const resetWorld = async (): Promise<T.CleanedTile[][]> => {
  const data: T.WorldResponse = await api.post("resetWorld", {});
  return data.message;
};

export const getEntityCognitiveMap = async (
  fullName: string
): Promise<T.CognitiveMapKnownTile[]> => {
  const query = new URLSearchParams({ fullName }).toString();
  const response = await fetch(address + `entityCognitiveMap?${query}`);
  const raw = await response.text();

  let data: unknown;
  try {
    data = JSON.parse(raw);
  } catch {
    if (raw.startsWith("Duplicate request detected")) {
      throw new Error(
        "Cognitive map endpoint unavailable (duplicate guard from default route). Restart backend."
      );
    }
    throw new Error("Cognitive map endpoint returned non-JSON response.");
  }

  if (!response.ok) {
    throw new Error(`Cognitive map request failed (${response.status}).`);
  }

  if (!data || typeof data !== "object" || !Array.isArray((data as { message?: unknown }).message)) {
    const message = (data as { message?: unknown })?.message;
    if (typeof message === "string" && message.includes("Welcome to the API")) {
      throw new Error("Cognitive map endpoint unavailable on backend. Restart backend.");
    }
    throw new Error("Unexpected cognitive map response shape.");
  }

  return (data as T.CognitiveMapResponse).message;
};

export default api;
