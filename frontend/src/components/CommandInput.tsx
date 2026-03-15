import React from "react";
import { Box, TextField, Button } from "@mui/material";

interface CommandInputProps {
  inputCommand: string;
  onChange: (value: string) => void;
  onKeyDown: (e: React.KeyboardEvent<HTMLTextAreaElement>) => void;
  onSend: () => void;
}

const CommandInput: React.FC<CommandInputProps> = ({
  inputCommand,
  onChange,
  onKeyDown,
  onSend,
}) => {
  return (
    <>
      <TextField
        multiline
        rows={3}
        value={inputCommand}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={(e) => onKeyDown(e as unknown as React.KeyboardEvent<HTMLTextAreaElement>)}
        placeholder="输入命令..."
        helperText="↵ 发送 | Shift+↵ 换行"
        fullWidth
        variant="outlined"
        slotProps={{
          htmlInput: {
            autoComplete: "off",
          },
        }}
      />
      <Box sx={{ mt: 1, display: "flex", justifyContent: "flex-end" }}>
        <Button sx={{ px: 5 }} variant="contained" onClick={onSend}>
          发送
        </Button>
      </Box>
    </>
  );
};

export default CommandInput;
