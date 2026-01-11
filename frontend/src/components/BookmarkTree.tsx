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
  Typography,
  alpha,
  Tooltip,
} from "@mui/material";
import {
  ExpandMore,
  ChevronRight,
  Edit as EditIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  FolderOutlined,
  BookmarkBorderOutlined,
} from "@mui/icons-material";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";

interface BookmarkGroup {
  name: string;
  bookmarks: SSHBookmark[];
}

interface BookmarkTreeProps {
  bookmarks: BookmarkGroup[];
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
  const [hoveredGroup, setHoveredGroup] = useState<string | null>(null);
  const [hoveredBookmark, setHoveredBookmark] = useState<string | null>(null);

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
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
      {/* 标题栏 */}
      <Box
        sx={{
          p: 2,
          pb: 1.5,
          borderBottom: 1,
          borderColor: "divider",
        }}
      >
        <Typography variant="h6" sx={{ fontWeight: 600, fontSize: "1.1rem" }}>
          书签
        </Typography>
      </Box>

      {/* 树形列表 */}
      <Box sx={{ flex: 1, overflowY: "auto" }}>
        <List sx={{ py: 1 }}>
          {bookmarks.map((group) => (
            <React.Fragment key={group.name}>
              {/* 分组 */}
              <ListItem
                disablePadding
                onMouseEnter={() => setHoveredGroup(group.name)}
                onMouseLeave={() => setHoveredGroup(null)}
                sx={{
                  "&:hover": {
                    backgroundColor: (theme) =>
                      alpha(theme.palette.primary.main, 0.04),
                  },
                }}
              >
                <ListItemButton
                  onClick={() => toggleGroup(group.name)}
                  sx={{
                    py: 0.75,
                    px: 1.5,
                  }}
                >
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      width: "100%",
                      gap: 1,
                    }}
                  >
                    {expandedGroups[group.name] ? (
                      <ExpandMore sx={{ fontSize: 20, color: "text.secondary" }} />
                    ) : (
                      <ChevronRight sx={{ fontSize: 20, color: "text.secondary" }} />
                    )}
                    <FolderOutlined sx={{ fontSize: 18, color: "warning.main" }} />
                    {editingGroup === group.name ? (
                      <TextField
                        value={newGroupName}
                        onChange={(e) => setNewGroupName(e.target.value)}
                        onKeyDown={(e) => handleGroupRenameKeyDown(e, group.name)}
                        onBlur={() => finishEditingGroup(group.name)}
                        autoFocus
                        size="small"
                        variant="standard"
                        sx={{ flex: 1, minWidth: 0 }}
                        inputProps={{
                          style: { fontSize: "0.95rem", padding: "2px 0" },
                        }}
                      />
                    ) : (
                      <ListItemText
                        primary={group.name}
                        primaryTypographyProps={{
                          fontSize: "0.95rem",
                          fontWeight: 500,
                        }}
                        sx={{ my: 0, flex: 1, minWidth: 0 }}
                      />
                    )}
                    <Box 
                      sx={{ 
                        display: "flex", 
                        gap: 0.5, 
                        ml: "auto",
                        visibility: hoveredGroup === group.name && editingGroup !== group.name ? "visible" : "hidden"
                      }}
                    >
                      <Tooltip title="编辑分组">
                        <IconButton
                          size="small"
                          onClick={(e) => {
                            e.stopPropagation();
                            startEditingGroup(group.name, group.name);
                          }}
                          sx={{ padding: "4px" }}
                        >
                          <EditIcon sx={{ fontSize: 16 }} />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="删除分组">
                        <IconButton
                          size="small"
                          onClick={(e) => {
                            e.stopPropagation();
                            onGroupDelete(group.name);
                          }}
                          sx={{ padding: "4px" }}
                        >
                          <DeleteIcon sx={{ fontSize: 16 }} />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </Box>
                </ListItemButton>
              </ListItem>

              {/* 书签列表 */}
              <Collapse in={expandedGroups[group.name]} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                  {group.bookmarks.map((bookmark) => (
                    <ListItem
                      key={bookmark.id}
                      disablePadding
                      onMouseEnter={() => setHoveredBookmark(bookmark.id)}
                      onMouseLeave={() => setHoveredBookmark(null)}
                    >
                      <ListItemButton
                        selected={selectedBookmark?.id === bookmark.id}
                        onClick={() => onBookmarkSelect(bookmark)}
                        sx={{
                          pl: 5,
                          py: 0.5,
                          "&.Mui-selected": {
                            backgroundColor: (theme) =>
                              alpha(theme.palette.primary.main, 0.12),
                            "&:hover": {
                              backgroundColor: (theme) =>
                                alpha(theme.palette.primary.main, 0.2),
                            },
                          },
                        }}
                      >
                        <BookmarkBorderOutlined
                          sx={{ fontSize: 16, mr: 1.5, color: "primary.main", flexShrink: 0 }}
                        />
                        <ListItemText
                          primary={bookmark.title}
                          secondary={`${bookmark.user}@${bookmark.host}`}
                          primaryTypographyProps={{
                            fontSize: "0.9rem",
                            noWrap: true,
                          }}
                          secondaryTypographyProps={{
                            fontSize: "0.75rem",
                            noWrap: true,
                          }}
                          sx={{ flex: 1, minWidth: 0 }}
                        />
                        <Tooltip title="删除书签">
                          <IconButton
                            size="small"
                            onClick={(e) => {
                              e.stopPropagation();
                              onBookmarkDelete(bookmark.id);
                            }}
                            sx={{ 
                              padding: "4px",
                              visibility: hoveredBookmark === bookmark.id ? "visible" : "hidden"
                            }}
                          >
                            <DeleteIcon sx={{ fontSize: 16 }} />
                          </IconButton>
                        </Tooltip>
                      </ListItemButton>
                    </ListItem>
                  ))}
                  {/* 添加书签按钮 */}
                  <ListItemButton
                    sx={{
                      pl: 5,
                      py: 0.5,
                      color: "primary.main",
                    }}
                    onClick={() => onBookmarkAdd(group.name)}
                  >
                    <AddIcon sx={{ fontSize: 16, mr: 1.5 }} />
                    <ListItemText
                      primary="添加书签"
                      primaryTypographyProps={{
                        fontSize: "0.9rem",
                      }}
                    />
                  </ListItemButton>
                </List>
              </Collapse>
            </React.Fragment>
          ))}

          {/* 添加分组 */}
          {newGroupInput ? (
            <ListItem sx={{ px: 1.5, py: 0.5 }}>
              <TextField
                value={newGroupNameInput}
                onChange={(e) => setNewGroupNameInput(e.target.value)}
                onKeyDown={handleNewGroupKeyDown}
                onBlur={handleAddNewGroup}
                autoFocus
                size="small"
                placeholder="分组名称"
                fullWidth
                sx={{
                  "& .MuiInputBase-root": {
                    fontSize: "0.9rem",
                  },
                }}
              />
            </ListItem>
          ) : (
            <ListItemButton
              onClick={() => setNewGroupInput(true)}
              sx={{
                py: 0.75,
                px: 1.5,
                color: "primary.main",
              }}
            >
              <AddIcon sx={{ fontSize: 20, mr: 1 }} />
              <ListItemText
                primary="添加分组"
                primaryTypographyProps={{
                  fontSize: "0.95rem",
                  fontWeight: 500,
                }}
              />
            </ListItemButton>
          )}
        </List>
      </Box>
    </Box>
  );
};

export default BookmarkTree;
