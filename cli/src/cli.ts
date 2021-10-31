import { program } from "commander";
import { authCommand } from "./auth/cmd";
import { createCommand } from "./create/cmd";

program.version("1.0");

program.addCommand(authCommand);
program.addCommand(createCommand);

program.parse(process.argv);
