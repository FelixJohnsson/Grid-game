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

interface PersonResponse {
    message: T.Person[];
    status: number;
}

export const getPersons = async (): Promise<T.Person[]> => {
    const data: PersonResponse = await api.get("people");
    return data.message;
};

interface WorldStateResponse {
    message: T.Building[];
    status: number;
}

export const getBuildings = async (): Promise<T.Building[]> => {
    const data: WorldStateResponse = await api.get("buildings");
    return data.message;
};

export default api;
