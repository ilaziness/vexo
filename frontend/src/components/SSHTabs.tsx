import { Box, Menu, MenuItem } from "@mui/material";
import React, { useRef, useCallback } from "react";
import SSHContainer from "./SSHContainer.tsx";
import { Events } from "@wailsio/runtime";
import { useSSHTabsStore, useReloadSSHTabStore } from "../stores/ssh.ts";
import { useTransferStore } from "../stores/transfer.ts";
import MainOpBar from "./MainOpBar.tsx";
import { DraggableTab } from "./DraggableTab.tsx";
import { useSSHContextMenu } from "../hooks/useSSHContextMenu.ts";
import { genTabIndex } from "../func/service.ts";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models.ts";
import {
  LogService,
  BookmarkService,
} from "../../bindings/github.com/ilaziness/vexo/services/index.ts";
import {
  DragDropContext,
  Droppable,
  DropResult,
} from "@hello-pangea/dnd";

const tabHeight = "40px";

export default function SSHTabs() {
  const {
    sshTabs,
    currentTab,
    delTab,
    pushTab,
    getByIndex,
    setCurrentTab,
    reorderTabs,
  } = useSSHTabsStore();
  const { doTabReload } = useReloadSSHTabStore();
  const { addProgress } = useTransferStore();
  const { anchorEl, tabIndex, isOpen, openMenu, closeMenu } = useSSHContextMenu();
  const scrollContainerRef = useRef<HTMLDivElement>(null);

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
    const number = sshTabs.length + 1;
    const newIndex = genTabIndex();
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
      sshInfo: currentTabInfo ? currentTabInfo.sshInfo : undefined,
    });
    setCurrentTab(newIndex);
  }, [tabIndex, sshTabs.length, getByIndex, pushTab, setCurrentTab, closeMenu]);

  const handleRefreshTab = useCallback(() => {
    closeMenu();
    if (tabIndex) {
      doTabReload(tabIndex);
    }
  }, [tabIndex, doTabReload, closeMenu]);

  const handleDragEnd = useCallback((result: DropResult) => {
    if (!result.destination) return;

    const sourceIndex = result.source.index;
    const destinationIndex = result.destination.index;

    if (sourceIndex === destinationIndex) return;

    reorderTabs(sourceIndex, destinationIndex);
  }, [reorderTabs]);

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
      <Box
        sx={{
          borderBottom: 1,
          borderColor: "divider",
          display: "flex",
          alignItems: "center",
          width: "100%",
          height: tabHeight,
          pr: 1,
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
                {sshTabs.map((item, index) => (
                  <DraggableTab
                    key={item.index}
                    item={item}
                    index={index}
                    isActive={currentTab === item.index}
                    onClick={() => setCurrentTab(item.index)}
                    onClose={(e) => {
                      e.stopPropagation();
                      delTab(item.index, currentTab);
                    }}
                    onContextMenu={(e) =>
                      openMenu(e as unknown as React.MouseEvent<HTMLElement>, item.index)
                    }
                  />
                ))}
                {provided.placeholder}
              </Box>
            )}
          </Droppable>
        </DragDropContext>
        <MainOpBar />
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
      <Box
        sx={{
          width: "100%",
          height: `calc(100% - ${tabHeight})`,
          position: "relative",
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
            <SSHContainer tabIndex={item.index} />
          </Box>
        ))}
      </Box>
    </Box>
  );
}
