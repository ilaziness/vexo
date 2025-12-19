import { Box, IconButton, Stack, Tooltip } from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import NoteAddIcon from "@mui/icons-material/NoteAdd";

export default function Header() {
  return (
    <Box
      component={"header"}
      sx={{
        background:
          "linear-gradient(to bottom, rgba(16, 18, 34, 1), rgba(27, 35, 63, 1))",
      }}
    >
      <Stack direction="column" spacing={1} alignItems="center" padding={0.5}>
        <Tooltip title="新建连接">
          <IconButton size="small">
            <NoteAddIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="设置">
          <IconButton size="small">
            <SettingsIcon />
          </IconButton>
        </Tooltip>
      </Stack>
    </Box>
  );
}
