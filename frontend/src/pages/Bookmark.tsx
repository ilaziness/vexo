import React, { useEffect, useState } from "react";
import { Box, Paper } from "@mui/material";
import {
  BookmarkService,
  LogService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";
import BookmarkTree from "../components/BookmarkTree";
import BookmarkForm from "../components/BookmarkForm";
import { useMessageStore } from "../stores/message";
import Message from "../components/Message";
import { parseCallServiceError } from "../func/service";
import { generateRandomId } from "../func/id";
import PasswordInputDialog from "../components/PasswordInputDialog";

interface BookmarkGroup {
  name: string;
  bookmarks: SSHBookmark[];
}

const Bookmark: React.FC = () => {
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
        // 过滤掉 null 值，确保类型安全
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
      // 查找旧分组
      const oldGroup = bookmarks.find((g) => g.name === oldGroupName);
      const bookmarksToMove = oldGroup?.bookmarks || [];

      // 为这些书签创建新ID并更新分组
      const updatedBookmarks = bookmarksToMove.map((b) => ({
        ...b,
        id: `${newGroupName}-${b.title}-${Date.now()}`, // 创建新ID以确保唯一性
        group_name: newGroupName,
      }));

      // 删除旧分组
      await BookmarkService.DeleteGroup(oldGroupName);

      // 添加更新后的书签到新分组
      for (const bookmark of updatedBookmarks) {
        await BookmarkService.SaveBookmark(bookmark);
      }

      // 重新加载书签
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
      if (selectedBookmark && selectedBookmark.group_name === groupName) {
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
      id: generateRandomId(),
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
      await SSHService.TestConnectInfo(
        bookmark.host,
        bookmark.port,
        bookmark.user,
        await BookmarkService.DecryptPassword(bookmark.password),
        bookmark.private_key,
        await BookmarkService.DecryptPassword(bookmark.private_key_password),
      );
      successMessage("连接测试成功");
    } catch (error) {
      LogService.Warn(`Connection test failed: ${error}`);
      errorMessage("连接测试失败");
    }
  };

  const handleSaveAndConnect = async (bookmark: SSHBookmark) => {
    try {
      // 先测试连接
      await SSHService.TestConnectInfo(
        bookmark.host,
        bookmark.port,
        bookmark.user,
        await BookmarkService.DecryptPassword(bookmark.password),
        bookmark.private_key,
        await BookmarkService.DecryptPassword(bookmark.private_key_password),
      );

      // 连接成功后保存书签
      await BookmarkService.SaveBookmark(bookmark);
      successMessage("书签保存成功");
      await loadBookmarks();

      // 调用ConnectBookmark触发事件，让SSHTabs处理创建tab
      await BookmarkService.ConnectBookmark(bookmark.id);
      BookmarkService.CloseWindow();
    } catch (error) {
      LogService.Warn(`Failed to save and connect: ${error}`);
      errorMessage("连接失败");
    }
  };

  return (
    <>
      <Box sx={{ height: "100%", display: "flex", overflow: "hidden" }}>
        {/* 左侧树形列表 */}
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

        {/* 右侧表单区域 */}
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
      <Message />
      <PasswordInputDialog />
    </>
  );
};

export default Bookmark;
