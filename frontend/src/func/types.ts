export interface FileInfo {
  name: string;
  size: number;
  mode: string;
  modeBits: number;
  uid: number;
  gid: number;
  modTime: string;
  isDir: boolean;
}

/** Build SFTP transfer owner key from bookmark or host identity. */
export function sftpOwnerKey(info: {
  bookmarkID?: string;
  user?: string;
  host?: string;
  port?: number;
}): string {
  if (info.bookmarkID) {
    return info.bookmarkID;
  }
  const user = info.user || "";
  const host = info.host || "";
  const port = info.port || 22;
  return `${user}@${host}:${port}`;
}

/** modeBits is Unix mode (0o7777), including setuid/setgid/sticky. */
export function modeBitsToOctal(modeBits: number): string {
  return (modeBits & 0o7777).toString(8);
}

export function parseOctalMode(input: string): number | null {
  const trimmed = input.trim().replace(/^0o/i, "");
  if (!/^[0-7]{3,4}$/.test(trimmed)) {
    return null;
  }
  return parseInt(trimmed, 8);
}
