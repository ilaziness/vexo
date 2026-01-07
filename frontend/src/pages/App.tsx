import { Box } from "@mui/material";
import SSHList from "../components/SSHList.tsx";
import Header from "../components/Header.tsx";
import Message from "../components/Message.tsx";

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
        <SSHList />
      </Box>
      <Message />
    </>
  );
}

export default App;
