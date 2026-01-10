import { Box, Tooltip, Typography } from "@mui/material";
import React from "react";

interface SSHTabTextProps {
  text: string;
}

const SSHTabText: React.FC<SSHTabTextProps> = ({ text }) => {
  return (
    <Box
      sx={{
        maxWidth: 130,
      }}
    >
      <Tooltip title={text}>
        <Typography noWrap>{text}</Typography>
      </Tooltip>
    </Box>
  );
};

export default SSHTabText;
