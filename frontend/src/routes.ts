import { createHashRouter } from "react-router";
import App from "./pages/App.tsx";
import Setting from "./pages/Setting";
import Command from "./pages/Command";

export default createHashRouter([
  {
    path: "/",
    Component: App,
  },
  {
    path: "/setting",
    Component: Setting,
  },
  {
    path: "/command",
    Component: Command,
  },
]);
