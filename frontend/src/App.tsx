import { Box, Snackbar } from '@mui/material';
import SSHList from './components/SSHList';
import Header from './components/Header';
import {useMessageStore} from './stores/common'

function App() {
  const { message,setClose } = useMessageStore();
  return (
    <>
      <Box sx={{ display: "flex", flexDirection: "row", height: "100%" }}>
        <Header />
        <SSHList />
      </Box>
      <Snackbar
        open={message.open}
        onClose={setClose}
        message={message.text}
        autoHideDuration={3000}
      />
    </>
  );
}

export default App
