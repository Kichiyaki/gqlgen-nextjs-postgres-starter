import { createContext } from "react";
import plTranslation from "./translations/pl";

const ctx = createContext(plTranslation);
ctx.displayName = "Translation";

export default ctx;
