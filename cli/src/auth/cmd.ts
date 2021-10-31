import { Command } from "commander";
import { login } from "./lib";

export const authCommand = new Command("login");

authCommand.action(async () => {
  await login();
});
