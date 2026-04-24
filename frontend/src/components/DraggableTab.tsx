import { Box, Tooltip } from "@mui/material";
import React from "react";
import CloseIcon from "@mui/icons-material/Close";
import { Draggable } from "@hello-pangea/dnd";
import { ConnectionStatus } from "../types/ssh";

const tabHeight = "40px";

export interface DraggableTabProps {
  item: {
    index: string;
    name: string;
    connectionStatus?: ConnectionStatus;
  };
  index: number;
  isActive: boolean;
  onClick: () => void;
  onClose: (e: React.MouseEvent) => void;
  onContextMenu: (e: React.MouseEvent) => void;
}

export const DraggableTab = React.memo(function DraggableTab({
  item,
  index,
  isActive,
  onClick,
  onClose,
  onContextMenu,
}: DraggableTabProps) {
  const getBackgroundColor = (isActive: boolean, isDragging: boolean) => {
    if (isActive) return "action.selected";
    if (isDragging) return "action.hover";
    return "inherit";
  };

  return (
    <Draggable draggableId={item.index} index={index}>
      {(provided, snapshot) => (
        <Box
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
          onClick={onClick}
          onContextMenu={onContextMenu}
          sx={{
            display: "flex",
            alignItems: "center",
            height: tabHeight,
            minHeight: tabHeight,
            minWidth: 150,
            width: 150,
            px: 1,
            gap: 0.5,
            borderRight: 2,
            borderColor: "divider",
            backgroundColor: getBackgroundColor(isActive, snapshot.isDragging),
            opacity: snapshot.isDragging ? 0.9 : 1,
            boxShadow: snapshot.isDragging
              ? "0 4px 8px rgba(0,0,0,0.3)"
              : "none",
            cursor: "pointer",
            userSelect: "none",
            "--wails-draggable": "no-drag",
            transition: "background-color 0.2s",
            "&:hover": {
              backgroundColor: isActive ? "action.selected" : "action.hover",
            },
            "&:active": {
              cursor: "grabbing",
            },
          }}
        >
          <Box
            sx={{
              width: 8,
              height: 8,
              borderRadius: "50%",
              flexShrink: 0,
              backgroundColor:
                item.connectionStatus === "connected"
                  ? "success.main"
                  : item.connectionStatus === "disconnected"
                    ? "error.main"
                    : "text.disabled",
            }}
          />
          <Tooltip title={item.name}>
            <Box
              sx={{
                flex: 1,
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
                fontSize: "0.875rem",
              }}
            >
              {item.name}
            </Box>
          </Tooltip>
          <Box
            component="span"
            onClick={onClose}
            sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              width: 20,
              height: 20,
              borderRadius: "50%",
              cursor: "pointer",
              opacity: 0.7,
              flexShrink: 0,
              "&:hover": {
                opacity: 1,
                backgroundColor: "action.hover",
              },
            }}
          >
            <CloseIcon sx={{ fontSize: 14 }} />
          </Box>
        </Box>
      )}
    </Draggable>
  );
});
