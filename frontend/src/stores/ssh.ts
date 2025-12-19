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
  pushTab: (tab: SSHTab) => void;
  delTab: (delTab: string, currentTab: string) => string; // 返回需要激活的标签
}

export const useSSHTabsStore = create<SSHTabs>((set, get) => ({
  sshTabs: [
    {
      index: `${getTabIndex()}`,
      name: "新建连接",
    },
  ] as SSHTab[],
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
    set({ sshTabs: newSshTabs });
    return newTabValue;
  },
}));
