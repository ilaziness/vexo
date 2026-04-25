import { Box } from "@mui/material";
import SSHTabs from "../components/SSHTabs.tsx";
import Header from "../components/subwindow/Header.tsx";
import Message from "../components/Message.tsx";
import HostKeyPrompt from "../components/HostKeyPrompt";

function SubMainWindow() {
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
        <SSHTabs />
      </Box>
      <Message />
      <HostKeyPrompt />
    </>
  );
}

export default SubMainWindow;
