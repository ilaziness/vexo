import { Snackbar } from "@mui/material";
import { useMessageStore } from "../stores/common.ts";

function Message() {
  const { message, setClose } = useMessageStore();
  return (
    <>
      <Snackbar
        open={message.open}
        onClose={setClose}
        message={message.text}
        autoHideDuration={3000}
      />
    </>
  );
}

export default Message;
