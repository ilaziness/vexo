import { createHashRouter } from "react-router";
import App from "./pages/App.tsx";
import Setting from "./pages/Setting";
import Bookmark from "./pages/Bookmark";

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
    path: "/bookmark",
    Component: Bookmark,
  },
]);
