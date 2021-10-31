import axios from "axios";
import { config } from "../config/lib";

export const create = async (url: string, api: string) => {
  let resp = await axios.post(
    api,
    { url },
    { headers: { Authorization: `Bearer ${config.get("token")}` } }
  );

  return resp.data.url;
};
