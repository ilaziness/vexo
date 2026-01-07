import React, { useState } from "react";
import {
  Box,
  List,
  ListItem,
  ListItemText,
  ListItemButton,
  Collapse,
  IconButton,
  TextField,
  InputAdornment,
} from "@mui/material";
import {
  ExpandMore,
  ExpandLess,
  Edit as EditIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
} from "@mui/icons-material";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";

interface BookmarkTreeProps {
  bookmarks: Record<string, SSHBookmark[]>;
  selectedBookmark: SSHBookmark | null;
  onBookmarkSelect: (bookmark: SSHBookmark) => void;
  onGroupRename: (oldGroupName: string, newGroupName: string) => void;
  onGroupAdd: (groupName: string) => void;
  onGroupDelete: (groupName: string) => void;
  onBookmarkAdd: (groupName: string) => void;
  onBookmarkDelete: (bookmarkId: string) => void;
}

interface GroupState {
  [key: string]: boolean;
}

const BookmarkTree: React.FC<BookmarkTreeProps> = ({
  bookmarks,
  selectedBookmark,
  onBookmarkSelect,
  onGroupRename,
  onGroupAdd,
  onGroupDelete,
  onBookmarkAdd,
  onBookmarkDelete,
}) => {
  const [expandedGroups, setExpandedGroups] = useState<GroupState>({});
  const [editingGroup, setEditingGroup] = useState<string | null>(null);
  const [newGroupName, setNewGroupName] = useState<string>("");
  const [newGroupInput, setNewGroupInput] = useState<boolean>(false);
  const [newGroupNameInput, setNewGroupNameInput] = useState<string>("");

  const toggleGroup = (groupName: string) => {
    setExpandedGroups((prev) => ({
      ...prev,
      [groupName]: !prev[groupName],
    }));
  };

  const startEditingGroup = (groupName: string, currentName: string) => {
    setEditingGroup(groupName);
    setNewGroupName(currentName);
  };

  const finishEditingGroup = (groupName: string) => {
    if (newGroupName.trim() && newGroupName !== groupName) {
      onGroupRename(groupName, newGroupName.trim());
    }
    setEditingGroup(null);
  };

  const handleGroupRenameKeyDown = (
    e: React.KeyboardEvent,
    groupName: string,
  ) => {
    if (e.key === "Enter") {
      finishEditingGroup(groupName);
    } else if (e.key === "Escape") {
      setEditingGroup(null);
    }
  };

  const handleAddNewGroup = () => {
    if (newGroupNameInput.trim()) {
      onGroupAdd(newGroupNameInput.trim());
      setNewGroupInput(false);
      setNewGroupNameInput("");
    }
  };

  const handleNewGroupKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleAddNewGroup();
    } else if (e.key === "Escape") {
      setNewGroupInput(false);
      setNewGroupNameInput("");
    }
  };

  return (
    <Box>
      <List>
        {Object.entries(bookmarks).map(([groupName, bookmarkList]) => (
          <React.Fragment key={groupName}>
            <ListItem disablePadding>
              <ListItemButton onClick={() => toggleGroup(groupName)}>
                {expandedGroups[groupName] ? <ExpandLess /> : <ExpandMore />}
                <ListItemText
                  primary={
                    editingGroup === groupName ? (
                      <TextField
                        value={newGroupName}
                        onChange={(e) => setNewGroupName(e.target.value)}
                        onKeyDown={(e) =>
                          handleGroupRenameKeyDown(e, groupName)
                        }
                        onBlur={() => setEditingGroup(null)}
                        autoFocus
                        size="small"
                        variant="standard"
                        inputProps={{ style: { fontSize: "1rem" } }}
                      />
                    ) : (
                      groupName
                    )
                  }
                  sx={{ ml: 1 }}
                />
                <IconButton
                  size="small"
                  onClick={(e) => {
                    e.stopPropagation();
                    startEditingGroup(groupName, groupName);
                  }}
                >
                  <EditIcon fontSize="small" />
                </IconButton>
                <IconButton
                  size="small"
                  onClick={(e) => {
                    e.stopPropagation();
                    onGroupDelete(groupName);
                  }}
                >
                  <DeleteIcon fontSize="small" />
                </IconButton>
              </ListItemButton>
            </ListItem>
            <Collapse
              in={expandedGroups[groupName]}
              timeout="auto"
              unmountOnExit
            >
              <List component="div" disablePadding>
                {bookmarkList.map((bookmark) => (
                  <ListItemButton
                    key={bookmark.id}
                    sx={{ pl: 4 }}
                    selected={selectedBookmark?.id === bookmark.id}
                    onClick={() => onBookmarkSelect(bookmark)}
                  >
                    <ListItemText primary={bookmark.title} />
                    <IconButton
                      size="small"
                      onClick={(e) => {
                        e.stopPropagation();
                        onBookmarkDelete(bookmark.id);
                      }}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </ListItemButton>
                ))}
                <ListItemButton
                  sx={{ pl: 4 }}
                  onClick={() => onBookmarkAdd(groupName)}
                >
                  <ListItemText
                    primary="+ 添加书签"
                    sx={{ color: "primary.main" }}
                  />
                </ListItemButton>
              </List>
            </Collapse>
          </React.Fragment>
        ))}
        {newGroupInput ? (
          <ListItem>
            <TextField
              value={newGroupNameInput}
              onChange={(e) => setNewGroupNameInput(e.target.value)}
              onKeyDown={handleNewGroupKeyDown}
              onBlur={() => setNewGroupInput(false)}
              autoFocus
              size="small"
              variant="standard"
              placeholder="分组名称"
              fullWidth
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={handleAddNewGroup} size="small">
                      <AddIcon fontSize="small" />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </ListItem>
        ) : (
          <ListItemButton onClick={() => setNewGroupInput(true)}>
            <ListItemText primary="+ 添加分组" sx={{ color: "primary.main" }} />
          </ListItemButton>
        )}
      </List>
    </Box>
  );
};

export default BookmarkTree;
