import { createHashRouter } from "react-router";
import App from "./pages/App.tsx";
import Setting from "./pages/Setting";

export default createHashRouter([
  {
    path: "/",
    Component: App,
  },
  {
    path: "/setting",
    Component: Setting,
  },
]);
