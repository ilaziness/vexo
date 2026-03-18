import React from "react";
import { Autocomplete, TextField, Checkbox, InputAdornment, IconButton } from "@mui/material";
import { Refresh } from "@mui/icons-material";
import { SSHTunnelSession } from "../types/command";

interface SessionSelectorProps {
  sessions: SSHTunnelSession[];
  selectedSessions: SSHTunnelSession[];
  onChange: (newValue: SSHTunnelSession[]) => void;
  onRefresh?: () => void; // 可选的刷新回调
}

const SessionSelector: React.FC<SessionSelectorProps> = ({
  sessions,
  selectedSessions,
  onChange,
  onRefresh,
}) => {
  return (
    <Autocomplete
      multiple
      options={sessions}
      getOptionLabel={(option) => option.clientKey}
      value={selectedSessions}
      onChange={(_, newValue) => onChange(newValue)}
      renderInput={(params) => (
        <TextField
          {...params}
          label="选择要发送到的 SSH 会话"
          placeholder="选择会话"
          size="small"
          slotProps={{
            input: {
              ...params.InputProps,
              endAdornment: (
                <InputAdornment position="end">
                  {onRefresh && (
                    <IconButton
                      onClick={onRefresh}
                      size="small"
                      edge="end"
                      title="刷新会话列表"
                    >
                      <Refresh fontSize="small" />
                    </IconButton>
                  )}
                </InputAdornment>
              ),
            },
          }}
        />
      )}
      renderOption={(props, option) => (
        <li {...props}>
          <Checkbox checked={selectedSessions.includes(option)} />
          {option.clientKey}
        </li>
      )}
      sx={{ my: 1 }}
    />
  );
};

export default SessionSelector;
