import { createHashRouter } from "react-router";
import App from "./pages/App.tsx";
import Setting from "./pages/Setting";
import Command from "./pages/Command";
import Tools from "./pages/Tools";
import ToolDetail from "./pages/ToolDetail";

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
  {
    path: "/tools",
    Component: Tools,
  },
  {
    path: "/tools/:toolId",
    Component: ToolDetail,
  },
]);
