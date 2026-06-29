import { useSSHTabsStore } from '../stores/ssh';
import { SSHContext } from '../../bindings/github.com/ilaziness/vexo/services/models';

export function getCurrentSSHContext(): SSHContext | undefined {
  const { currentTab, getByIndex } = useSSHTabsStore.getState();
  const tab = getByIndex(currentTab);
  const sshInfo = tab?.sshInfo;

  if (!sshInfo?.linkID || !sshInfo.host || !sshInfo.user || typeof sshInfo.port !== 'number' || sshInfo.port <= 0) {
    return undefined;
  }

  return new SSHContext({
    link_id: sshInfo.linkID,
    host: sshInfo.host,
    port: sshInfo.port,
    user: sshInfo.user,
  });
}
