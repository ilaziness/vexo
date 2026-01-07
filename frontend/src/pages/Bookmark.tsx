import React, { useEffect, useState } from "react";
import { Box, Paper, Divider } from "@mui/material";
import { BookmarkService } from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";
import BookmarkTree from "../components/BookmarkTree";
import BookmarkForm from "../components/BookmarkForm";
import { useMessageStore } from "../stores/common";
import Message from "../components/Message";
import { parseCallServiceError } from "../func/service";

const Bookmark: React.FC = () => {
  const [bookmarks, setBookmarks] = useState<Record<string, SSHBookmark[]>>({});
  const [selectedBookmark, setSelectedBookmark] = useState<SSHBookmark | null>(
    null,
  );

  const { errorMessage } = useMessageStore();

  useEffect(() => {
    loadBookmarks();
  }, []);

  const loadBookmarks = async () => {
    try {
      const config = await BookmarkService.ListBookmarks();
      if (config && Array.isArray(config)) {
        // 将BookmarkGroup数组转换为按组名索引的对象
        const bookmarkMap: Record<string, SSHBookmark[]> = {};
        for (const group of config) {
          if (group && group.name && Array.isArray(group.bookmarks)) {
            bookmarkMap[group.name] = group.bookmarks;
          }
        }
        setBookmarks(bookmarkMap);
      } else {
        setBookmarks({});
      }
    } catch (error) {
      console.error("Failed to load bookmarks:", error);
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
      // 删除旧分组的所有书签
      const bookmarksToMove = bookmarks[oldGroupName] || [];

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
        await BookmarkService.AddBookmark(bookmark);
      }

      // 重新加载书签
      await loadBookmarks();
    } catch (error) {
      console.error("Failed to rename group:", error);
      errorMessage("重命名分组失败: " + parseCallServiceError(error));
    }
  };

  const handleGroupAdd = async (groupName: string) => {
    try {
      await BookmarkService.AddGroup(groupName);
      await loadBookmarks();
    } catch (error) {
      console.error("Failed to add group:", error);
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
      console.error("Failed to delete group:", error);
      errorMessage("删除分组失败: " + parseCallServiceError(error));
    }
  };

  const handleBookmarkAdd = async (groupName: string) => {
    try {
      const newBookmark: SSHBookmark = {
        id: `${groupName}-new-${Date.now()}`,
        title: `新书签 ${bookmarks[groupName]?.length ? bookmarks[groupName].length + 1 : 1}`,
        group_name: groupName,
        host: "",
        port: 22,
        private_key: "",
        user: "",
        password: "",
      };
      await BookmarkService.AddBookmark(newBookmark);
      setSelectedBookmark(newBookmark);
      await loadBookmarks();
    } catch (error) {
      console.error("Failed to add bookmark:", error);
      errorMessage("添加书签失败: " + parseCallServiceError(error));
    }
  };

  const handleBookmarkDelete = async (bookmarkId: string) => {
    try {
      await BookmarkService.DeleteBookmark(bookmarkId);
      await loadBookmarks();
      if (selectedBookmark?.id === bookmarkId) {
        setSelectedBookmark(null);
      }
    } catch (error) {
      console.error("Failed to delete bookmark:", error);
      errorMessage("删除书签失败: " + parseCallServiceError(error));
    }
  };

  const handleSaveBookmark = async (bookmark: SSHBookmark) => {
    try {
      await BookmarkService.AddBookmark(bookmark);
      await loadBookmarks();
    } catch (error) {
      console.error("Failed to save bookmark:", error);
      errorMessage("保存书签失败: " + parseCallServiceError(error));
    }
  };

  const handleTestConnection = (bookmark: SSHBookmark) => {
    // TODO: 实现测试连接功能
    console.log("Testing connection for:", bookmark);
  };

  const handleSaveAndConnect = async (bookmark: SSHBookmark) => {
    try {
      await BookmarkService.AddBookmark(bookmark);
      await loadBookmarks();
      // TODO: 实现连接功能
      console.log("Saving and connecting to:", bookmark);
    } catch (error) {
      console.error("Failed to save and connect:", error);
      errorMessage("保存并连接失败: " + parseCallServiceError(error));
    }
  };

  return (
    <>
      <Box sx={{ height: "100vh", display: "flex", flexDirection: "column" }}>
        <Box sx={{ display: "flex", flex: 1, overflow: "hidden" }}>
          <Paper
            sx={{
              width: "30%",
              minWidth: 300,
              height: "100%",
              overflow: "auto",
            }}
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
          <Divider orientation="vertical" flexItem />
          <Paper
            sx={{
              flex: 1,
              height: "100%",
              overflow: "auto",
            }}
          >
            <BookmarkForm
              bookmark={selectedBookmark}
              onSave={handleSaveBookmark}
              onTestConnection={handleTestConnection}
              onSaveAndConnect={handleSaveAndConnect}
            />
          </Paper>
        </Box>
      </Box>
      <Message />
    </>
  );
};

export default Bookmark;
