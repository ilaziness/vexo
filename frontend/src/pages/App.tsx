import { Box } from "@mui/material";
import SSHTabs from "../components/SSHTabs.tsx";
import Header from "../components/Header.tsx";
import Message from "../components/Message.tsx";
import PasswordInputDialog from "../components/PasswordInputDialog.tsx";
import HostKeyPrompt from "../components/HostKeyPrompt";
import { AppService } from "../../bindings/github.com/ilaziness/vexo/services";

function App() {
  return (
    <>
      <Box
        sx={{
          display: "flex",
          flexDirection: "row",
          height: "100%",
          width: "100%",
        }}
      >
        <Header />
        <SSHTabs
          onMinimize={() => AppService.MainWindowMin()}
          onMaximize={() => AppService.MainWindowMax()}
          onClose={() => AppService.MainWindowClose()}
        />
      </Box>
      <Message />
      <PasswordInputDialog />
      <HostKeyPrompt />
    </>
  );
}

export default App;
