import { DeviceFlowClient } from "oauth_device_flow";
import { config } from "../config/lib";

const client = new DeviceFlowClient(
  {
    audience: "https://lnk.wie.gg",
    scopes: ["create:link"],
    client_id: "Cq4ugV3lY0JczABgY4th64XvPQRInsJ0",
    code_url: "https://wiegg.eu.auth0.com/oauth/device/code",
    token_url: "https://wiegg.eu.auth0.com/oauth/token",
  },
  {
    cache: {
      beforeCacheAccess: () => {
        if (config.has("token_cache")) {
          return JSON.parse("token_cache");
        }
      },
      afterCacheAccess: (cache) => {
        config.set("token_cache", JSON.stringify(cache));
      },
    },
  }
);

export const login = async () => {
  const token = await client.acquireToken();

  config.set("token", token);

  return token;
};

export const authMiddleware = async () => {
  const token = await client.acquireTokenSilently();

  config.set("token", token);
};
