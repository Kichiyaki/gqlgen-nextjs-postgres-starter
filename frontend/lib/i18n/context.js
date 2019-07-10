import { createContext } from "react";
import plTranslation from "./languages/pl";

const ctx = createContext(plTranslation);
ctx.displayName = "Translation";

export default ctx;
