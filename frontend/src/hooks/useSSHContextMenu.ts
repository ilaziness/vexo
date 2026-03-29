import { useState, useCallback } from "react";

export interface ContextMenuState {
  anchorEl: HTMLElement | null;
  tabIndex: string | null;
}

export function useSSHContextMenu() {
  const [state, setState] = useState<ContextMenuState>({
    anchorEl: null,
    tabIndex: null,
  });

  const openMenu = useCallback((e: React.MouseEvent, tabIndex: string) => {
    e.preventDefault();
    setState({
      anchorEl: e.currentTarget as HTMLElement,
      tabIndex,
    });
  }, []);

  const closeMenu = useCallback(() => {
    setState({ anchorEl: null, tabIndex: null });
  }, []);

  return {
    anchorEl: state.anchorEl,
    tabIndex: state.tabIndex,
    isOpen: Boolean(state.anchorEl),
    openMenu,
    closeMenu,
  };
}
