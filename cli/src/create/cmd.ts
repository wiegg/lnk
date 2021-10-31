import { Command } from "commander";
import { create } from "./lib";

export const createCommand = new Command("create")
  .alias("c")
  .argument("<url>")
  .option("--api <api>", "API path", "http://lnk.wie.gg")
  .action(async (url, opts) => {
    const shorturl = await create(url, opts.api);

    console.log(`Your url was succesfully shortened: ${opts.api}${shorturl}`);
  });
