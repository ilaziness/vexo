import { Box, Menu, MenuItem } from "@mui/material";
import React, { useRef, useCallback } from "react";
import SSHTabBody from "./SSHTabBody.tsx";
import { Events } from "@wailsio/runtime";
import { useSSHTabsStore, useReloadSSHTabStore } from "../stores/ssh";
import { ConnectionStatus } from "../types/ssh";
import { useTransferStore } from "../stores/transfer.ts";
import { useAIAssistantStore } from "../stores/aiAssistant";
import OpBar from "./OpBar.tsx";
import { DraggableTab } from "./DraggableTab.tsx";
import { useSSHContextMenu } from "../hooks/useSSHContextMenu";
import { genTabIndex } from "../func/service";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models";
import {
  LogService,
  BookmarkService,
} from "../../bindings/github.com/ilaziness/vexo/services/index.ts";
import { DragDropContext, Droppable, DropResult } from "@hello-pangea/dnd";
import AISideBar from "./ai/AISideBar";

interface SSHTabsProps {
  onClose?: () => void;
}

// tab栏的高度
const tabHeight = "40px";

export default function SSHTabs({ onClose }: SSHTabsProps) {
  const sshTabs = useSSHTabsStore((state) => state.sshTabs);
  const currentTab = useSSHTabsStore((state) => state.currentTab);
  const delTab = useSSHTabsStore((state) => state.delTab);
  const pushTab = useSSHTabsStore((state) => state.pushTab);
  const getByIndex = useSSHTabsStore((state) => state.getByIndex);
  const setCurrentTab = useSSHTabsStore((state) => state.setCurrentTab);
  const reorderTabs = useSSHTabsStore((state) => state.reorderTabs);
  const doTabReload = useReloadSSHTabStore((state) => state.doTabReload);
  const addProgress = useTransferStore((state) => state.addProgress);
  const { anchorEl, tabIndex, isOpen, openMenu, closeMenu } =
    useSSHContextMenu();
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const sidebarOpen = useAIAssistantStore((state) => state.sidebarOpen);
  const sidebarWidth = useAIAssistantStore((state) => state.sidebarWidth);

  React.useEffect(() => {
    const unsubscribeProgress = Events.On("eventProgress", (event: any) => {
      LogService.Debug(`SSHTabs eventProgress: ${JSON.stringify(event)}`);
      const eventData = event.data as ProgressData;
      addProgress(eventData);
    });

    const unsubscribeConnectBookmark = Events.On(
      "eventConnectBookmark",
      async (event: any) => {
        try {
          const bookmarkID = event.data as string;
          LogService.Debug(`Connecting to bookmark: ${bookmarkID}`);

          const bookmark =
            await BookmarkService.GetBookmarkForConnect(bookmarkID);
          if (!bookmark) {
            LogService.Warn(`Bookmark not found: ${bookmarkID}`);
            return;
          }

          const newTab = {
            index: genTabIndex(),
            name: bookmark.title,
            sshInfo: {
              bookmarkID: bookmarkID,
              host: bookmark.host,
              port: bookmark.port,
              user: bookmark.user,
              proxyJumpID: bookmark.proxy_jump_id,
            },
            connectionStatus: ConnectionStatus.Connecting,
          };

          pushTab(newTab);
          setCurrentTab(newTab.index);
        } catch (error) {
          LogService.Warn(`Failed to connect bookmark: ${error}`);
        }
      },
    );

    return () => {
      unsubscribeProgress();
      unsubscribeConnectBookmark();
    };
  }, [addProgress, pushTab, setCurrentTab]);

  const handleCloseTab = useCallback(() => {
    if (tabIndex === null) return;
    delTab(tabIndex, currentTab);
    closeMenu();
  }, [tabIndex, currentTab, delTab, closeMenu]);

  const handleDuplicateTab = useCallback(() => {
    closeMenu();
    const currentTabInfo = getByIndex(tabIndex || "");
    const newIndex = genTabIndex();
    pushTab({
      index: newIndex,
      name: currentTabInfo?.name || '',
      sshInfo: currentTabInfo?.sshInfo
        ? { ...currentTabInfo.sshInfo, linkID: undefined }
        : undefined,
      connectionStatus: currentTabInfo?.sshInfo
        ? ConnectionStatus.Connecting
        : undefined,
    });
    setCurrentTab(newIndex);
  }, [tabIndex, sshTabs.length, getByIndex, pushTab, setCurrentTab, closeMenu]);

  const handleRefreshTab = useCallback(() => {
    closeMenu();
    if (tabIndex) {
      doTabReload(tabIndex);
    }
  }, [tabIndex, doTabReload, closeMenu]);

  const handleDragEnd = useCallback(
    (result: DropResult) => {
      if (!result.destination) return;

      const sourceIndex = result.source.index;
      const destinationIndex = result.destination.index;

      if (sourceIndex === destinationIndex) return;

      reorderTabs(sourceIndex, destinationIndex);
    },
    [reorderTabs],
  );

  return (
    <Box
      sx={{
        height: "100%",
        width: "calc(100% - 42px)",
        display: "flex",
        flexDirection: "column",
        flex: 1,
      }}
    >
      {/* tabs */}
      <Box
        sx={{
          borderBottom: 1,
          borderColor: "divider",
          display: "flex",
          alignItems: "center",
          width: "100%",
          height: tabHeight,
          pr: 1,
          bgcolor: "background.paper",
          "--wails-draggable": "drag",
        }}
      >
        <DragDropContext onDragEnd={handleDragEnd}>
          <Droppable droppableId="ssh-tabs" direction="horizontal">
            {(provided) => (
              <Box
                ref={(el: HTMLDivElement | null) => {
                  provided.innerRef(el);
                  scrollContainerRef.current = el;
                }}
                {...provided.droppableProps}
                sx={{
                  display: "flex",
                  flex: 1,
                  height: tabHeight,
                  overflow: "auto",
                  scrollbarWidth: "none",
                  "&::-webkit-scrollbar": {
                    display: "none",
                  },
                }}
              >
                {sshTabs.map((tab, index) => (
                  <DraggableTab
                    key={tab.index}
                    item={tab}
                    index={index}
                    isActive={currentTab === tab.index}
                    onClick={() => setCurrentTab(tab.index)}
                    onClose={(e) => {
                      e.stopPropagation();
                      delTab(tab.index, currentTab);
                    }}
                    onContextMenu={(e) =>
                      openMenu(
                        e as unknown as React.MouseEvent<HTMLElement>,
                        tab.index,
                      )
                    }
                  />
                ))}
                {provided.placeholder}
              </Box>
            )}
          </Droppable>
        </DragDropContext>
        <OpBar draggable={false} showBorder={false} onClose={onClose} />
        <Menu
          open={isOpen}
          anchorEl={anchorEl}
          onClose={closeMenu}
          anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
          transformOrigin={{ vertical: "top", horizontal: "center" }}
        >
          <MenuItem onClick={handleCloseTab}>关闭</MenuItem>
          <MenuItem onClick={handleDuplicateTab}>复制</MenuItem>
          <MenuItem onClick={handleRefreshTab}>刷新</MenuItem>
        </Menu>
      </Box>

      {/* content */}
      <Box
        sx={{
          width: "100%",
          height: `calc(100% - ${tabHeight})`,
          position: "relative",
        }}
      >
        <Box
          sx={{
            width: "100%",
            height: "100%",
            position: "relative",
            boxSizing: "border-box",
            paddingRight: sidebarOpen ? `${sidebarWidth}px` : 0,
          }}
        >
          {sshTabs.map((item) => (
            <Box
              role="tabpanel"
              key={item.index}
              sx={{
                display: "flex",
                width: "100%",
                height: "100%",
                position: "absolute",
                top: 0,
                left: currentTab === item.index ? 0 : "-9999rem",
              }}
            >
              <SSHTabBody tabIndex={item.index} />
            </Box>
          ))}
        </Box>
        <AISideBar />
      </Box>
    </Box>
  );
}
