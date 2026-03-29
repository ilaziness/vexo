import { Box } from "@mui/material";
import React from "react";
import TerminalIcon from "@mui/icons-material/Terminal";
import CloseIcon from "@mui/icons-material/Close";
import { Draggable } from "@hello-pangea/dnd";

const tabHeight = "40px";

export interface DraggableTabProps {
  item: {
    index: string;
    name: string;
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
            backgroundColor: isActive
              ? "action.selected"
              : snapshot.isDragging
                ? "action.hover"
                : "inherit",
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
          }}
        >
          <TerminalIcon sx={{ fontSize: 14, flexShrink: 0 }} />
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
