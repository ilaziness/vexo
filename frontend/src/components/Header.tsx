import { Box, IconButton, Stack, Tooltip } from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import NoteAddIcon from "@mui/icons-material/NoteAdd";
import BookmarksIcon from "@mui/icons-material/Bookmarks";
import {
  BookmarkService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { useSSHTabsStore } from "../stores/ssh";
import { getTabIndex } from "../func/service";

export default function Header() {
  const { pushTab } = useSSHTabsStore();

  const showSettingWindow = () => {
    ConfigService.ShowWindow().then(() => {});
  };

  const handleAddTab = () => {
    const number = useSSHTabsStore.getState().sshTabs.length + 1;
    const newIndex = `${getTabIndex()}`;
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
    });
  };

  return (
    <Box
      component={"header"}
      sx={{
        width: "42px",
        background:
          "linear-gradient(to bottom, rgba(16, 18, 34, 1), rgba(27, 35, 63, 1))",
      }}
    >
      <Stack direction="column" spacing={1} alignItems="center" padding={0.5}>
        <Tooltip title="新建连接">
          <IconButton size="small" onClick={handleAddTab}>
            <NoteAddIcon />
          </IconButton>
        </Tooltip>
        <Tooltip
          title="书签"
          onClick={() => {
            BookmarkService.ShowWindow().then(() => {});
          }}
        >
          <IconButton size="small">
            <BookmarksIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="设置">
          <IconButton size="small" onClick={showSettingWindow}>
            <SettingsIcon />
          </IconButton>
        </Tooltip>
      </Stack>
    </Box>
  );
}
