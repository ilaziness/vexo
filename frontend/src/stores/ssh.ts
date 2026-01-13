import { create } from "zustand";
import { SSHLinkInfo, SSHTab } from "../types/ssh";
import { getTabIndex } from "../func/service";

export const useSSHLinksStore = create((set) => ({
  sshLinks: [] as SSHLinkInfo[],
  setSSHLink: (link: SSHLinkInfo) =>
    set((state) => ({
      sshLinks: [...state.sshLinks, link],
    })),
}));

export interface SSHTabs {
  sshTabs: SSHTab[];
  currentTab: string;
  pushTab: (tab: SSHTab) => void;
  delTab: (delTab: string, currentTab: string) => string; // 返回需要激活的标签
  setName: (index: string, name: string) => void;
  setCurrentTab: (index: string) => void;
  getByIndex: (index: string) => SSHTab | undefined;
  setSSHInfo: (index: string, sshInfo: SSHLinkInfo) => void;
}

export const useSSHTabsStore = create<SSHTabs>((set, get) => ({
  sshTabs: [
    {
      index: `${getTabIndex()}`,
      name: "新建连接",
    },
  ] as SSHTab[],
  currentTab: `${getTabIndex()}`,
  pushTab: (tab: SSHTab) =>
    set((state: any) => ({
      sshTabs: [...state.sshTabs, tab],
    })),
  delTab: (delTab: string, currentTab: string): string => {
    const tabIndex = get().sshTabs.findIndex((item) => item.index === delTab);
    if (tabIndex === -1) {
      return "0";
    }
    // 删除指定索引的标签
    let newSshTabs = get().sshTabs.filter((item) => item.index !== delTab);

    // 如果删除的是当前激活的标签，需要切换到另一个标签
    let newTabValue = currentTab;
    if (currentTab === delTab) {
      if (newSshTabs.length > 0) {
        newTabValue = newSshTabs[newSshTabs.length - 1].index;
      } else {
        // 如果没有其他标签了，创建一个新标签
        newSshTabs = [
          {
            index: `${getTabIndex()}`,
            name: "新建连接",
          },
        ];
        newTabValue = newSshTabs[0].index;
      }
    }
    set({ sshTabs: newSshTabs, currentTab: newTabValue });
    return newTabValue;
  },
  setName: (index: string, name: string) =>
    set((state) => ({
      sshTabs: state.sshTabs.map((tab) =>
        tab.index === index ? { ...tab, name } : tab,
      ),
    })),
  setCurrentTab: (index: string) =>
    set({ currentTab: index }),
  getByIndex: (index: string): SSHTab | undefined =>
    get().sshTabs.find((val) => {
      return val.index === index;
    }),
  setSSHInfo: (index: string, sshInfo: SSHLinkInfo) =>
    set((state) => ({
      sshTabs: state.sshTabs.map((tab) =>
        tab.index === index ? { ...tab, sshInfo } : tab,
      ),
    })),
}));

export interface reloadSSHTabStoreType {
  reloadTab: {
    index: string;
    count: number;
  };
  doTabReload: (index: string) => void;
}

export const useReloadSSHTabStore = create<reloadSSHTabStoreType>(
  (set, get) => ({
    reloadTab: {
      index: "",
      count: 0,
    },
    doTabReload: (index: string) => {
      set({ reloadTab: { index: index, count: get().reloadTab.count } });
    },
  }),
);
