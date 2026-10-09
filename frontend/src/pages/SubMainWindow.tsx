import { Box } from "@mui/material";
import SSHTabs from "../components/SSHTabs.tsx";
import Header from "../components/subwindow/Header.tsx";
import Message from "../components/Message.tsx";
import HostKeyPrompt from "../components/HostKeyPrompt";
import KeyboardInteractivePrompt from "../components/KeyboardInteractivePrompt";
import PasswordInputDialog from "../components/PasswordInputDialog";
import ShortcutListener from "../components/ShortcutListener";

function SubMainWindow() {
  return (
    <>
      <ShortcutListener />
      <Box
        sx={{
          display: "flex",
          flexDirection: "row",
          height: "100%",
          width: "100%",
        }}
      >
        <Header />
        <SSHTabs />
      </Box>
      <Message />
      <HostKeyPrompt />
      <KeyboardInteractivePrompt />
      <PasswordInputDialog />
    </>
  );
}

export default SubMainWindow;
