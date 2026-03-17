import React, { useEffect, useState } from "react";
import { Box, Paper } from "@mui/material";
import {
  BookmarkService,
  LogService,
  SSHBookmark,
} from "../../bindings/github.com/ilaziness/vexo/services";
import BookmarkTree from "./BookmarkTree";
import BookmarkForm from "./BookmarkForm";
import { useMessageStore } from "../stores/message";
import { parseCallServiceError } from "../func/service";

interface BookmarkGroup {
  name: string;
  bookmarks: SSHBookmark[];
}

interface BookmarkProps {
  onRequestClose?: () => void;
}

const Bookmark: React.FC<BookmarkProps> = ({ onRequestClose }) => {
  const [bookmarks, setBookmarks] = useState<BookmarkGroup[]>([]);
  const [selectedBookmark, setSelectedBookmark] = useState<SSHBookmark | null>(
    null,
  );
  const { errorMessage, successMessage } = useMessageStore();

  useEffect(() => {
    loadBookmarks();
  }, []);

  const loadBookmarks = async () => {
    try {
      const config = await BookmarkService.ListBookmarks();
      if (config && Array.isArray(config)) {
        const validBookmarks = config.filter(
          (group): group is BookmarkGroup =>
            group !== null && group !== undefined,
        );
        setBookmarks(validBookmarks);
      } else {
        setBookmarks([]);
      }
    } catch (error) {
      LogService.Warn(`Failed to load bookmarks:${error}`);
      errorMessage("加载书签失败: " + parseCallServiceError(error));
    }
  };

  const handleBookmarkSelect = (bookmark: SSHBookmark) => {
    setSelectedBookmark(bookmark);
  };

  const handleGroupRename = async (
    oldGroupName: string,
    newGroupName: string,
  ) => {
    try {
      await BookmarkService.UpdateGroup(oldGroupName, newGroupName);
      await loadBookmarks();
    } catch (error) {
      LogService.Warn(`Failed to rename group: ${error}`);
      errorMessage("重命名分组失败: " + parseCallServiceError(error));
    }
  };

  const handleGroupAdd = async (groupName: string) => {
    try {
      await BookmarkService.AddGroup(groupName);
      await loadBookmarks();
    } catch (error) {
      LogService.Warn(`Failed to add group: ${error}`);
      errorMessage("添加分组失败: " + parseCallServiceError(error));
    }
  };

  const handleGroupDelete = async (groupName: string) => {
    try {
      await BookmarkService.DeleteGroup(groupName);
      await loadBookmarks();
      if (selectedBookmark?.group_name === groupName) {
        setSelectedBookmark(null);
      }
    } catch (error) {
      LogService.Warn(`Failed to delete group: ${error}`);
      errorMessage("删除分组失败: " + parseCallServiceError(error));
    }
  };

  const handleBookmarkAdd = (groupName: string) => {
    const group = bookmarks.find((g) => g.name === groupName);
    const bookmarkCount = group?.bookmarks?.length || 0;
    const newBookmark: SSHBookmark = {
      id: "",
      title: `新书签 ${bookmarkCount + 1}`,
      group_name: groupName,
      host: "",
      port: 22,
      private_key: "",
      private_key_password: "",
      user: "",
      password: "",
    };
    setSelectedBookmark(newBookmark);
  };

  const handleBookmarkDelete = async (bookmarkId: string) => {
    try {
      await BookmarkService.DeleteBookmark(bookmarkId);
      await loadBookmarks();
      if (selectedBookmark?.id === bookmarkId) {
        setSelectedBookmark(null);
      }
    } catch (error) {
      LogService.Warn(`Failed to delete bookmark: ${error}`);
      errorMessage("删除书签失败: " + parseCallServiceError(error));
    }
  };

  const handleSaveBookmark = async (bookmark: SSHBookmark) => {
    try {
      await BookmarkService.SaveBookmark(bookmark);
      successMessage("书签保存成功");
      await loadBookmarks();
    } catch (error) {
      LogService.Warn(`Failed to save bookmark: ${error}`);
      errorMessage("保存书签失败: " + parseCallServiceError(error));
    }
  };

  const handleTestConnection = async (bookmark: SSHBookmark) => {
    try {
      await BookmarkService.TestConnection(bookmark);
      successMessage("连接测试成功");
    } catch (error) {
      LogService.Warn(`Connection test failed: ${error}`);
      errorMessage("连接测试失败: " + parseCallServiceError(error));
    }
  };

  const handleSaveAndConnect = async (bookmark: SSHBookmark) => {
    try {
      await BookmarkService.SaveAndConnect(bookmark);
      successMessage("书签保存并连接成功");
      await loadBookmarks();
      onRequestClose?.();
    } catch (error) {
      LogService.Warn(`Failed to save and connect: ${error}`);
      errorMessage("保存并连接失败: " + parseCallServiceError(error));
    }
  };

  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        overflow: "hidden",
      }}
    >
      <Paper
        sx={{
          width: 280,
          flexShrink: 0,
          borderRight: 1,
          borderColor: "divider",
          display: "flex",
          flexDirection: "column",
        }}
        elevation={0}
        square
      >
        <BookmarkTree
          bookmarks={bookmarks}
          selectedBookmark={selectedBookmark}
          onBookmarkSelect={handleBookmarkSelect}
          onGroupRename={handleGroupRename}
          onGroupAdd={handleGroupAdd}
          onGroupDelete={handleGroupDelete}
          onBookmarkAdd={handleBookmarkAdd}
          onBookmarkDelete={handleBookmarkDelete}
        />
      </Paper>
      <Box
        sx={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          overflow: "hidden",
        }}
      >
        <BookmarkForm
          bookmark={selectedBookmark}
          groupNames={bookmarks.map((g) => g.name)}
          onSave={handleSaveBookmark}
          onTestConnection={handleTestConnection}
          onSaveAndConnect={handleSaveAndConnect}
        />
      </Box>
    </Box>
  );
};

export default Bookmark;
