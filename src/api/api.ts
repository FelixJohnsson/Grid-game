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

export default api;
